package app

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"Player/internal/AudioEngine"
	"Player/internal/lyrics"
	"Player/internal/media"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Song struct {
	metadata media.Metadata
	lyrics   *lyrics.Lyrics
}

func (s Song) Title() string       { return s.metadata.Title }
func (s Song) Description() string { return s.metadata.Artist }
func (s Song) FilterValue() string { return s.metadata.Title }

type playerState int

const (
	stateStopped playerState = iota
	statePlaying
	statePaused
)

type tickMsg time.Time
type loadingTickMsg time.Time
type songsLoadedMsg struct {
	songs    []Song
	musicDir string
}
type lyricsLoadedMsg struct {
	song   *Song
	lyrics *lyrics.Lyrics
}

type model struct {
	songs          []Song
	list           list.Model
	progress       progress.Model
	state          playerState
	currentSong    *Song
	currentTime    float64
	engine         *AudioEngine.FFplayEngine
	mu             sync.Mutex
	width          int
	height         int
	lastUpdateTime time.Time
	shuffle        bool
	playHistory    []int
	lyricsLoading  bool
	loading        bool
	loadingDots    int
	seeking        bool
	musicDir       string
}

type keyMap struct {
	Play     key.Binding
	Pause    key.Binding
	Stop     key.Binding
	Next     key.Binding
	Previous key.Binding
	Forward  key.Binding
	Backward key.Binding
	Shuffle  key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Play: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "play"),
	),
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "pause/resume"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	Next: key.NewBinding(
		key.WithKeys("n", "right"),
		key.WithHelp("n/→", "next"),
	),
	Previous: key.NewBinding(
		key.WithKeys("b", "left"),
		key.WithHelp("b/←", "previous"),
	),
	Forward: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "forward 5s"),
	),
	Backward: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rewind 5s"),
	),
	Shuffle: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "shuffle"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	// add vol up(c) and down(v)
}

func initialModel(musicDir string) *model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Playlist"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	prog := progress.New(progress.WithDefaultGradient())

	rand.Seed(time.Now().UnixNano())

	return &model{
		songs:          []Song{},
		list:           l,
		progress:       prog,
		state:          stateStopped,
		engine:         AudioEngine.NewFFplayEngine(),
		lastUpdateTime: time.Now(),
		shuffle:        false,
		playHistory:    make([]int, 0),
		lyricsLoading:  false,
		loading:        true,
		loadingDots:    0,
		musicDir:       musicDir,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(loadingTickCmd(), tea.EnterAltScreen)
}

func loadingTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		return loadingTickMsg(t)
	})
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		rightWidth := msg.Width / 3
		listHeight := msg.Height - 12
		if listHeight < 5 {
			listHeight = 5
		}

		m.list.SetSize(rightWidth-8, listHeight)

		return m, nil

	case songsLoadedMsg:
		m.songs = msg.songs
		m.musicDir = msg.musicDir
		items := make([]list.Item, len(m.songs))
		for i, song := range m.songs {
			items[i] = song
		}
		m.list.SetItems(items)
		m.loading = false

		if len(m.songs) > 0 {
			return m, m.playSongCmd(&m.songs[0])
		}
		return m, tickCmd()

	case loadingTickMsg:
		if m.loading {
			m.loadingDots = (m.loadingDots + 1) % 4
			return m, loadingTickCmd()
		}
		return m, tickCmd()

	case lyricsLoadedMsg:
		if m.currentSong != nil && m.currentSong.metadata.FilePath == msg.song.metadata.FilePath {
			m.currentSong.lyrics = msg.lyrics
			m.lyricsLoading = false
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			if key.Matches(msg, keys.Quit) {
				return m, tea.Quit
			}
			return m, nil
		}

		if key.Matches(msg, keys.Quit) {
			m.stopPlayback()
			return m, tea.Quit
		}

		if m.list.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, keys.Play):
			if len(m.songs) > 0 {
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					selectedSong := selectedItem.(Song)
					// Find the song in the main list to ensure we have the correct pointer/reference
					for i := range m.songs {
						if m.songs[i].metadata.FilePath == selectedSong.metadata.FilePath {
							return m, m.playSongCmd(&m.songs[i])
						}
					}
				}
			}
			return m, nil

		case key.Matches(msg, keys.Pause):
			if m.state == statePlaying {
				m.pausePlayback()
			} else if m.state == statePaused {
				m.resumePlayback()
			} else if m.state == stateStopped && len(m.songs) > 0 {
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					selectedSong := selectedItem.(Song)
					for i := range m.songs {
						if m.songs[i].metadata.FilePath == selectedSong.metadata.FilePath {
							return m, m.playSongCmd(&m.songs[i])
						}
					}
				}
			}
			return m, nil

		case key.Matches(msg, keys.Stop):
			m.stopPlayback()
			return m, nil

		case key.Matches(msg, keys.Next):
			return m, m.playNextCmd()

		case key.Matches(msg, keys.Previous):
			return m, m.playPreviousCmd()

		case key.Matches(msg, keys.Forward):
			if !m.seeking {
				m.seek(5)
			}
			return m, nil

		case key.Matches(msg, keys.Backward):
			if !m.seeking {
				m.seek(-5)
			}
			return m, nil

		case key.Matches(msg, keys.Shuffle):
			m.shuffle = !m.shuffle
			m.playHistory = make([]int, 0)
			return m, nil
		}

	case tickMsg:
		if m.state == statePlaying && m.currentSong != nil {
			now := time.Now()
			elapsed := now.Sub(m.lastUpdateTime).Seconds()

			if elapsed > 0 && elapsed < 1.0 {
				m.currentTime += elapsed
			}
			m.lastUpdateTime = now

			if m.currentTime >= m.currentSong.metadata.Duration {
				return m, tea.Batch(m.playNextCmd(), tickCmd())
			}
		}
		return m, tickCmd()
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) playSongCmd(song *Song) tea.Cmd {
	m.playSong(song)

	if song.lyrics == nil && !m.lyricsLoading {
		m.lyricsLoading = true
		return loadLyricsAsync(song, m.musicDir)
	}

	return nil
}

func (m *model) playSong(song *Song) {
	if song == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.engine.Stop()
	m.currentSong = song
	m.currentTime = 0
	m.lastUpdateTime = time.Now()
	m.state = statePlaying
	m.seeking = false
	m.lyricsLoading = false

	m.engine.SetOnComplete(func() {
		m.mu.Lock()
		if m.state == statePlaying {
			m.state = stateStopped
		}
		m.mu.Unlock()
	})

	m.engine.Play(song.metadata.FilePath, 0, 100)
}

func (m *model) stopPlayback() {
	m.engine.Stop()
	m.state = stateStopped
	m.currentTime = 0
	m.seeking = false
}

func (m *model) pausePlayback() {
	m.engine.Pause()
	m.state = statePaused
}

func (m *model) resumePlayback() {
	if m.currentSong != nil && m.state == statePaused {
		m.state = statePlaying
		m.lastUpdateTime = time.Now()
		m.seeking = false

		m.engine.SetOnComplete(func() {
			m.mu.Lock()
			if m.state == statePlaying {
				m.state = stateStopped
			}
			m.mu.Unlock()
		})

		m.engine.Resume(m.currentTime, 100)
	}
}

func (m *model) seek(seconds float64) {
	if m.currentSong == nil || m.seeking {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.seeking = true
	defer func() {
		time.AfterFunc(100*time.Millisecond, func() {
			m.seeking = false
		})
	}()

	newTime := m.currentTime + seconds
	if newTime < 0 {
		newTime = 0
	}
	if newTime >= m.currentSong.metadata.Duration {
		m.seeking = false
		m.playNextCmd()
		return
	}

	wasPlaying := m.state == statePlaying

	m.currentTime = newTime

	if wasPlaying {
		m.state = statePlaying
		m.lastUpdateTime = time.Now()

		m.engine.SetOnComplete(func() {
			m.mu.Lock()
			if m.state == statePlaying {
				m.state = stateStopped
			}
			m.mu.Unlock()
		})

		m.engine.Seek(m.currentTime, 100)
	}
}

func (m *model) playNextCmd() tea.Cmd {
	if len(m.songs) == 0 {
		return nil
	}

	var nextIdx int

	if m.shuffle {
		if len(m.playHistory) >= len(m.songs) {
			m.playHistory = make([]int, 0)
		}

		unplayed := make([]int, 0)
		for i := 0; i < len(m.songs); i++ {
			played := false
			for _, h := range m.playHistory {
				if h == i {
					played = true
					break
				}
			}
			if !played {
				unplayed = append(unplayed, i)
			}
		}

		if len(unplayed) > 0 {
			nextIdx = unplayed[rand.Intn(len(unplayed))]
		} else {
			nextIdx = rand.Intn(len(m.songs))
		}

		m.playHistory = append(m.playHistory, nextIdx)
	} else {
		currentIdx := m.list.Index()
		nextIdx = (currentIdx + 1) % len(m.songs)
	}

	m.list.Select(nextIdx)
	return m.playSongCmd(&m.songs[nextIdx])
}

func (m *model) playPreviousCmd() tea.Cmd {
	if len(m.songs) == 0 {
		return nil
	}

	var prevIdx int

	if m.shuffle && len(m.playHistory) > 1 {
		m.playHistory = m.playHistory[:len(m.playHistory)-1]
		prevIdx = m.playHistory[len(m.playHistory)-1]
		m.playHistory = m.playHistory[:len(m.playHistory)-1]
	} else {
		currentIdx := m.list.Index()
		prevIdx = currentIdx - 1
		if prevIdx < 0 {
			prevIdx = len(m.songs) - 1
		}
	}

	m.list.Select(prevIdx)
	return m.playSongCmd(&m.songs[prevIdx])
}

func (m *model) View() string {
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			Align(lipgloss.Center).
			Width(m.width).
			Height(m.height)

		dots := strings.Repeat(".", m.loadingDots)
		spaces := strings.Repeat(" ", 3-m.loadingDots)

		loadingText := fmt.Sprintf("\n\n♪ Music Player\n\nLoading songs%s%s\n\nPlease wait...", dots, spaces)

		return loadingStyle.Render(loadingText)
	}

	if m.width == 0 {
		return "Initializing..."
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("63")).
		Padding(0, 1)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	var leftPanel string
	if m.currentSong != nil {
		meta := m.currentSong.metadata

		leftPanel = fmt.Sprintf("%s\n\n", titleStyle.Render("♪ Now Playing"))
		leftPanel += fmt.Sprintf("Title:  %s\n", meta.Title)
		leftPanel += fmt.Sprintf("Artist: %s\n", meta.Artist)
		leftPanel += fmt.Sprintf("Album:  %s\n\n", meta.Album)

		stateStr := "■ Stopped"
		switch m.state {
		case statePlaying:
			stateStr = "▶ Playing"
		case statePaused:
			stateStr = "❚❚ Paused"
		}
		leftPanel += fmt.Sprintf("Status: %s\n", stateStr)

		if m.shuffle {
			leftPanel += "Mode:   🔀 Shuffle ON\n\n"
		} else {
			leftPanel += "Mode:   ▶ Sequential\n\n"
		}

		progressPercent := 0.0
		if meta.Duration > 0 {
			progressPercent = m.currentTime / meta.Duration
			if progressPercent > 1.0 {
				progressPercent = 1.0
			}
		}

		leftPanel += m.progress.ViewAs(progressPercent) + "\n"

		currentMin := int(m.currentTime) / 60
		currentSec := int(m.currentTime) % 60
		totalMin := int(meta.Duration) / 60
		totalSec := int(meta.Duration) % 60
		leftPanel += fmt.Sprintf("%02d:%02d / %02d:%02d\n\n",
			currentMin, currentSec, totalMin, totalSec)
	} else {
		leftPanel = titleStyle.Render("♪ Music Player") + "\n\n"
		leftPanel += "No song playing\n"
		leftPanel += "Select a song and press 'p' to play\n"
		leftPanel += "or press 'space' to start\n\n"
	}

	leftPanel += infoStyle.Render(fmt.Sprintf("\nControls:\n" +
		"  p: play selected  space: pause/resume\n" +
		"  s: stop           n/→: next  b/←: prev\n" +
		"  t: forward 5s     r: rewind 5s\n" +
		"  h: shuffle        q: quit\n"))

	lyricsSection := "\n" + titleStyle.Render("Lyrics") + "\n\n"

	if m.currentSong != nil {
		if m.lyricsLoading {
			lyricsSection += infoStyle.Render("Loading lyrics...\n\n")
		} else if m.currentSong.lyrics != nil && m.currentSong.lyrics.Loaded {
			if len(m.currentSong.lyrics.Lines) == 0 {
				lyricsSection += infoStyle.Render("No lyrics available for this song.\n\n")
			} else {
				current, next := getCurrentLyrics(m.currentSong.lyrics, m.currentTime)

				currentStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("230")).
					Bold(true)

				nextStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("241"))

				if current != "" {
					lyricsSection += currentStyle.Render(current) + "\n"
				}
				if next != "" {
					lyricsSection += nextStyle.Render(next) + "\n"
				}
				lyricsSection += "\n"
			}
		} else {
			lyricsSection += infoStyle.Render("Lyrics not loaded\n\n")
		}
	} else {
		lyricsSection += infoStyle.Render("No song playing\n\n")
	}

	leftPanel += lyricsSection

	rightPanel := m.list.View()

	if m.currentSong != nil {
		meta := m.currentSong.metadata

		audioInfoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			PaddingTop(1).
			MarginTop(1)

		audioInfo := audioInfoStyle.Render(
			fmt.Sprintf("Bitrate:     %s\n", meta.Bitrate) +
				fmt.Sprintf("Codec:       %s\n", meta.Codec) +
				fmt.Sprintf("Sample Rate: %s", meta.SampleRate),
		)

		rightPanel += "\n" + audioInfo
	}

	leftWidth := (m.width * 2 / 3) - 4
	rightWidth := (m.width / 3) - 2

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Padding(1, 2)

	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Padding(1, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftPanel),
		rightStyle.Render(rightPanel),
	)
}

func loadLyricsAsync(song *Song, musicDir string) tea.Cmd {
	return func() tea.Msg {
		if lrc, found := lyrics.LoadFromFile(song.metadata.FilePath, musicDir); found {
			return lyricsLoadedMsg{song: song, lyrics: &lrc}
		}

		content, err := lyrics.FetchFromAPI(song.metadata.Artist, song.metadata.Title, song.metadata.Album)
		if err != nil {
			return lyricsLoadedMsg{song: song, lyrics: &lyrics.Lyrics{Loaded: true}}
		}

		if _, err := lyrics.SaveToFile(song.metadata.FilePath, musicDir, content); err != nil {
			return lyricsLoadedMsg{song: song, lyrics: &lyrics.Lyrics{Loaded: true}}
		}

		lines := lyrics.Parse(content)
		return lyricsLoadedMsg{song: song, lyrics: &lyrics.Lyrics{Lines: lines, Loaded: true}}
	}
}

func getCurrentLyrics(ly *lyrics.Lyrics, currentTime float64) (current, next string) {
	if ly == nil || !ly.Loaded || len(ly.Lines) == 0 {
		return "", ""
	}

	currentIdx := -1
	for i, line := range ly.Lines {
		if line.Time <= currentTime {
			currentIdx = i
		} else {
			break
		}
	}

	if currentIdx >= 0 {
		current = ly.Lines[currentIdx].Text
		if currentIdx+1 < len(ly.Lines) {
			next = ly.Lines[currentIdx+1].Text
		}
	} else if len(ly.Lines) > 0 {
		next = ly.Lines[0].Text
	}

	return current, next
}
