<h2
      font-family="Arial Black, Helvetica, sans-serif" 
      font-size="28" 
      font-weight="bold" 
      text-anchor="middle" 
      letterSpacing="2">
      R S S D U K E S T E R
</h2>


 - A simple clean terminal-based RSS feed reader. 
 - Read, organize, and search your favorite news feeds right from your terminal.
 - Application utilizes sqlite3, and has pre-subsrcibed feeds.
 
<details>
<summary> <h2>Features</h2> </summary>

<details>
<summary>Manage RSS feeds</summary>
<img src="images/home.png" alt="Image 1">

</details>


<details>
<summary>View posts from feeds</summary>

<img src="images/feed.png" alt="Image 2">
</details>


<details>
<summary>Read articles in terminal</summary>

 <img src="images/article.png" alt="Image 3">
</details>

<details>
<summary>Search across all your feeds</summary>

 <img src="images/search.png" alt="Image 4">
</details>

</details>

<details>
<summary> <h2> Build from source</h2> </summary>

<h3>requirements</h3>

- Go 1.21 or later
- GCC compiler
- SQLite3
- Windows Terminal or PowerShell (CMD not supported)

<details>
<summary> <h3>windows</h3> </summary>


1. **Install MSYS2** 
   ```powershell
   winget install MSYS2.MSYS2
   ```

2. **Open MSYS2 MINGW64** (from Start Menu) and run:
   ```bash
   pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-sqlite3
   ```

3. **Add MinGW to temporary PATH** (in PowerShell):
   ```powershell
   $env:Path += ";C:\msys64\mingw64\bin"
   ```

4. **Build**
   ```powershell
   git clone https://github.com/IvanYaremko/rssdukester.git
   cd rssdukester
   $env:CGO_ENABLED=1
   go build
   ```

5. **Run**
   ```powershell
   .\rssdukester.exe
   ```
</details>

<details>
<summary> <h3>linux</h3> </summary>

1. **Install dependencies**
   ```bash
   sudo apt-get update
   sudo apt-get install gcc libsqlite3-dev
   ```

2. **Build**
   ```bash
   git clone https://github.com/IvanYaremko/rssdukester.git
   cd rssdukester
   CGO_ENABLED=1 go build
   ```

3. **Run**
   ```bash
   ./rssdukester
   ```
</details>

<details>
<summary> <h3>macOS</h3> </summary>

1. **Install dependencies**
   ```bash
   brew install sqlite3
   ```

2. **Build**
   ```bash
   git clone https://github.com/IvanYaremko/rssdukester.git
   cd rssdukester
   CGO_ENABLED=1 go build
   ```

3. **Run**
   ```bash
   ./rssdukester
   ```
</details> 

<details>
<summary> <h3>troubleshooting</h3> </summary>

If you see `gcc: executable file not found in %PATH%`:
1. Make sure you opened MSYS2 MINGW64 and ran the pacman command
2. Verify GCC is installed by running: `gcc --version`
3. Ensure you added MinGW to PATH as shown in the build steps
</details>
</details>

<details>
<summary> <h2>Under the hood</h2> </summary>

`r s s d u k e s t e r` uses:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- [lipgloss](https://github.com/charmbracelet/lipgloss) for text formatting
- [go-readability](https://github.com/go-shiori/go-readability)
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) 
- [sqlite](https://www.sqlite.org/)

</details>
