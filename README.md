Message Notification Daemon
=====================

This is a desktop message notification daemon that can be left to run in the background to check multiple E-Mail accounts for new mail. POP3 is only supported for checking mail boxes and support for other protocols will be added in the future.

## Usage
```
Usage of mnd:
  -d="": The directory in which the database will be stored
  -f="": The configuration file in which the user settings are stored
  -g=false: Generate a configuration file
  -v=false: debugging/verbose information
```

## Requirements

	- GNU/Linux only
	- libnotify (GNU/Linux) for notifications

## Features

	- Supports multiple accounts

## TODO

	- Support proxys
	- Support Pushover
	- Support custom sounds
	- Support playing sounds if notification daemon can not play them

## Screenshot

The look of notifications will differ depending on your current window and icon theme. If your notification daemon (such as xfce4-notifyd) supports playing sounds, it will also play the "mail-unread" sound which will also differ depending on the current theme.

![ScreenShot](http://apollo.firebit.co.uk/~dc0/imgsrc/2015-03-14--1426373296_524x342_scrot.png)