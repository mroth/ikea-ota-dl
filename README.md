Go port of [ikea-ota-downloader.py][1], a simple script to download the most recent Ikea smart home accessory firmware images, in order to perform OTA updates via deCONZ et al for those who use their own Zigbee equipment and not a first-party IKEA gateway/software.

_Why?_ I tend to run this script on a very minimal server and I don't like to keep a copy of Python around just for that, so this way I can use a compiled binary instead.


    Usage: ikea-ota-dl [options] <destination>

    Options:
      -feed string
            firmware update feed URL (default "http://fw.ota.homesmart.ikea.net/feed/version_info.json")
      -v    verbose (default true)
      -workers int
            max concurrent downloads (default 4)

[1]: https://github.com/dresden-elektronik/deconz-rest-plugin/blob/master/ikea-ota-download.py
