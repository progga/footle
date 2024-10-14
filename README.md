Footle
======
Footle is a [debugger front-end](https://en.wikipedia.org/wiki/Debugger#Debugger_front-ends) for the [Xdebug](https://xdebug.org/) PHP debugger.  It offers a browser-based user interface for Xdebug.  The goal is to make interactive debugging easy for PHP newcomers.

## Development status
Footle is heavily under construction.  It works, but user experience is not good enough yet.  It cannot cope with multiple simultaneous debugging sessions.  It also lacks user authentication.

## Installation
- Open Footle's [release page](https://github.com/progga/footle/releases)
- Download the appropriate archive file for your platform.  e.g. footle-linux-64.tar.gz for 64-bit GNU/Linux environment.
- Uncompress the archive file.  e.g. `tar xfz footle-linux-64.tar.gz`, `unzip footle-win-64.zip`, etc.
- If you are using Xdebug **2.6** or prior, try an older [version of Footle](https://github.com/progga/footle/releases/tag/2018-11-26-dev).

## Usage
### Configuring Xdebug
For Footle to work, you first need to properly [install and configure Xdebug](https://xdebug.org/docs/install). Setting up Xdebug is out of scope of this document.  Instead we point you to the [Xdebug configuration wizard](https://xdebug.org/wizard.php).

To keep things simple, here is a sample xdebug.ini file that works for **simple** PHP projects:
```
zend_extension=xdebug.so

xdebug.mode = debug
```

Many of us use the `xdebug.discover_client_host` (formerly `xdebug.remote_connect_back`) or `xdebug.client_host` (formerly `xdebug.remote_host`) settings to instruct Xdebug to talk to Xdebug clients in a different machine.  This setting is usually not needed for Footle as it is easier to run Footle in the same machine as Xdebug.

If you suspect that Xdebug is not talking to Footle, launch Footle with the **-vvv** option.  This should dump into screen all traffic to and from Xdebug.  Total silence would mean Xdebug is not talking to Footle.  In which case review the xdebug settings present in the output of the phpinfo() function.

### Launching Footle
```
$ cd /var/www/html/
$ ~/footle-linux-64/footle
```

This makes two assumptions:
- Your PHP code is inside the /var/www/html/ directory.
- Footle is running in the **same** machine as Xdebug.

Press Ctrl-C to quit.

### Debugging process
- Open http://localhost:1234/ in a browser.  This brings up the Footle user interface.  When Footle is running in a different machine from the browser, use that machine's hostname rather than _localhost_.
- You should be presented with Footle's file picker.  The file picker is always the first tab *within* Footle's interface.  This should list all files and directories from your PHP codebase.
- Click one or more PHP files from the file picker. Selected files will open in their own tabs. Note that these tabs are not browser tabs. These tabs are part of the webpage drawn by Footle.
- Set breakpoints by clicking line numbers. Line numbers appear at the left edge of each file.  Alternatively, insert the `xdebug_break()` function within the PHP source file where a breakpoint is desired.
- Now in another browser tab or window, open a webpage that will execute the PHP files where you have just set breakpoints.
- Once execution reaches the breakpoint, the line with the breakpoint is highlighted by a light-green background.
- To inspect local and global variables, use the two buttons labelled *Locals* and *Globals*

## Supported platforms
Footle is cross-platform.  We prepare distributions for FreeBSD, GNU/Linux, MacOS, and Windows.  Minimum web browser requirement is [Firefox 60 ESR](https://en.wikipedia.org/wiki/History_of_Firefox#Rapid_release_with_ESR) or [Chromium](https://en.wikipedia.org/wiki/Chromium_(web_browser)) 69.  Recent browsers of other flavours may work although none are tested as yet.

Footle is not fully responsive yet, so at least a tablet display is recommended.

## Issues
Have a question about Footle?  Found a bug?  Please [file an issue](https://github.com/progga/footle/issues/new).
## Licence
[Simplified BSD licence](https://en.wikipedia.org/wiki/BSD_licenses#2-clause).  It effectively means you can do whatever you want with Footle.
