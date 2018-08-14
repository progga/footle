Footle
======
Footle is a [debugger front-end](https://en.wikipedia.org/wiki/Debugger#Debugger_front-ends) for the Xdebug PHP debugger.  It offers a browser-based user interface for Xdebug.  The goal is to make interactive debugging easy for PHP newcomers.

## Development status
Footle is heavily under construction.  It works, but user experience is not good enough yet.  Also, it cannot cope with multiple simultaneous debugging sessions.

## Installation
- Open Footle's [release page](https://github.com/progga/footle/releases)
- Download the appropriate archive file for your platform.  e.g. footle-linux-64.tar.gz for 64-bit GNU/Linux environment.
- Uncompress the archive file.  e.g. `tar xfz footle-linux-64.tar.gz`, `unzip footle-win-64.zip`, etc.

## Usage
### Configuring Xdebug
For Footle to work, you first need to properly install and configure Xdebug. Setting up Xdebug is out of scope for this document.  Instead we will point you to the [Xdebug configuration wizard](https://xdebug.org/wizard.php).

To keep things simple, we are providing a sample xdebug.ini file that works for **simple** PHP projects:
```
zend_extension=xdebug.so

xdebug.remote_enable = 1
xdebug.remote_autostart = 1
```

### Launching Footle
`$ ./footle/bin/footle -codebase /var/www/html`

This makes two assumptions:
- Your PHP code is inside the /var/www/html/ directory.
- Footle is running in the same machine as Xdebug.

### Debugging process
- Open http://localhost:1234/ in a browser.  This brings up the Footle user interface.
- You should be presented with Footle's file picker.  The file picker is always the first tab *within* Footle's interface.  This should list all files and directories from /var/www/html/
- Click one or more PHP files from the file picker.  Selected files will open in their own tabs.  Note that these tabs are *not* browser tabs.  These tabs are part of the webpage drawn by Footle.
- Set breakpoints by clicking on line numbers.  Line numbers appear at the left edge of each file.
- Now in another browser tab or window, open a webpage that will execute the PHP files where you have just set breakpoints.
- Once execution reaches the breakpoint, the line with the breakpoint is highlighted by a light-green background.
- To inspect local and global variables, use the two buttons labelled *Locals* and *Globals*

## Supported platforms
Footle is cross-platform.  We prepare distributions for FreeBSD, GNU/Linux, MacOS, and Windows.  Minimum web browser requirement is [Firefox 52 ESR](https://en.wikipedia.org/wiki/History_of_Firefox#Extended_Support_Release).  Recent browsers of other flavours may work although none are tested as yet.

## Licence
[Simplified BSD licence](https://en.wikipedia.org/wiki/BSD_licenses#2-clause).  It effectively means you can do whatever you want with Footle.
