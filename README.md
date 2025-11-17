# BG3 Watchdog

## Description
BG3 Watchdog is used in conjunction with running Baldur's Gate 3 and [Sunshine](https://github.com/LizardByte/Sunshine) streaming. Larian designed BG3 to start as `bg3.exe` but it opens a Larian Launcher. Running the game from the launcher causes the launcher to close, terminating the `bg3.exe` and opening a new `bg3.exe` to run the actual game. Sunshine streaming watches for when processes are terminated to send events to any streaming sessions, indicating the game has been exited. Instead of watching for `bg3.exe`, Sunshine can use this project's `bg3Watchdog.exe` to start the game and wait for it to be terminated. 

`bg3Watchdog.exe` will start `bg3.exe` and wait for the launcher to close, bridging the gap where the process will seem to close and reopen. The `bg3Watchdog.exe` will stay running until the actual game is exited, and then it will close, allowing Sunshine to have a better indicator of if the game is running.

## Building
run `go build`

## Prerequesites
* Steam
* Baldur's Gate 3
* Sunshine

## Usage
Place the built exe in a directory somewhere accessible for Sunshine, and select it as the executable for running Balder's Gate 3

## Config Defaults
Found in config.yaml in the directory with `bg3Watchdog.exe`

```yaml
steam_path: "C:\\Program Files (x86)\\Steam\\steam.exe"
bg3_app_id: 1086940
launcher_dwell_seconds: 10
main_app_poll_seconds: 2
debug: false
```