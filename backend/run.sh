#!/bin/sh
echo "Starting app..."
(binscrape &)
echo $(pidof binscrape)
(api &)
echo $(pidof api)
