!define APP_NAME "StellePlayer"
!define APP_VERSION "1.0.0"
!define APP_PUBLISHER "StellePlayer"
!define INSTALL_DIR "$PROGRAMFILES\${APP_NAME}"
!define EXE_NAME "StellePlayer.exe"
!define BUILD_DIR "..\build"
SetCompressor /SOLID lzma

Name "${APP_NAME} ${APP_VERSION}"
OutFile "${APP_NAME}_Setup_${APP_VERSION}.exe"
InstallDir "${INSTALL_DIR}"
RequestExecutionLevel admin

Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install"
  SetOutPath "$INSTDIR"

  File /r "..\build\${EXE_NAME}"

  ; Create CLI aliases using batch files
  FileOpen $0 "$INSTDIR\Splayer.bat" w
  FileWrite $0 '@echo off$\r$\n"%~dp0${EXE_NAME}" %*'
  FileClose $0

  FileOpen $0 "$INSTDIR\Splay.bat" w
  FileWrite $0 '@echo off$\r$\n"%~dp0${EXE_NAME}" %*'
  FileClose $0

  CreateDirectory "$SMPROGRAMS\${APP_NAME}"
  CreateShortCut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" "$INSTDIR\${EXE_NAME}"
  CreateShortCut "$DESKTOP\${APP_NAME}.lnk" "$INSTDIR\${EXE_NAME}"

  WriteUninstaller "$INSTDIR\Uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$INSTDIR\Uninstall.exe"

  ; Install and run PATH setup script
  SetOutPath "$INSTDIR\Scripts"
  File "add_stelle_player_to_path.bat"
  DetailPrint "Running PATH setup script..."
  ExecWait '"$INSTDIR\Scripts\add_stelle_player_to_path.bat"'
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\${EXE_NAME}"
  Delete "$INSTDIR\Uninstall.exe"
  Delete "$INSTDIR\Splayer.bat"
  Delete "$INSTDIR\Splay.bat"
  RMDir /r "$INSTDIR"

  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"
SectionEnd
