#!/bin/sh

/api &
exec nginx -g 'daemon off;'
