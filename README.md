# ttyTmux

ttyTmux is a Serial Setup and Logging Utility, intended as Daily Cron job in Lab. An opinionated wrapper around Tmux, it solves the problem of not just setting up Tmux sessions, but cycling logs daily to correlate events. One major flaw is that the baudrate must be the same for all devices, which is impractical in most environments but acceptable for our lab's focus of monitoring OLT. Since minicom must be run as root, tmux is run as root, and therefore must be accessed with "sudo tmux ls" etc after initialization.

| Flag | Description |
| ------ | ------ |
| -h | Show this help |
| -b | Baudate (for all devices, sorry! |

# Example Usage
```sh
~/test> ls /dev/ttyUSB*
/dev/ttyUSB0  /dev/ttyUSB1  /dev/ttyUSB2
~/test> tmux ls
no server running on /tmp/tmux-1000/default
~/test> sudo tmux ls
no server running on /tmp/tmux-0/default
~/test> sudo ps aux | grep minicom
lbnp       39701  0.0  0.0   9040   664 pts/0    S+   10:28   0:00 grep --color=auto minicom
~/test> sudo ./ttyTmux
~/test> sudo tmux ls
tty0: 1 windows (created Tue Sep 28 10:28:10 2021)
tty1: 1 windows (created Tue Sep 28 10:28:10 2021)
tty2: 1 windows (created Tue Sep 28 10:28:10 2021)
~/test> ls
20210928_tty0.log  20210928_tty1.log  20210928_tty2.log  ttyTmux*
~/test> sudo ps aux | grep minicom
root       39892  0.0  0.0   9848  4048 pts/1    S+   10:28   0:00 minicom -D /dev/ttyUSB0 -b 115200 -C 20210928_tty0.log
root       39919  0.0  0.0   9848  4060 pts/2    S+   10:28   0:00 minicom -D /dev/ttyUSB1 -b 115200 -C 20210928_tty1.log
root       39934  0.0  0.0   9848  4084 pts/3    S+   10:28   0:00 minicom -D /dev/ttyUSB2 -b 115200 -C 20210928_tty2.log
lbnp       39965  0.0  0.0   9040   736 pts/0    S+   10:28   0:00 grep --color=auto minicom
```
