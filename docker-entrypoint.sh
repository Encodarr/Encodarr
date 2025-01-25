#!/bin/ash
PUID=${PUID:-1000}
PGID=${PGID:-1000}

start_app() {
    if id -u $PUID > /dev/null 2>&1 && id -g $PGID > /dev/null 2>&1; then
        chown -R $PUID:$PGID /config /movies /series /transcode
        su -m $PUID -c 'setsid /app/transfigurr' &
        APP_PID=$!
    else
        setsid /app/transfigurr &
        APP_PID=$!
    fi
}

kill_app() {
    echo "Killing processes..."
    
    # Kill main process group
    if [ ! -z "$APP_PID" ]; then
        pkill -TERM -P $APP_PID
        kill -TERM -$APP_PID 2>/dev/null
    fi

    # Force kill any remaining ffmpeg processes
    pkill -9 ffmpeg 2>/dev/null
    
    # Wait for processes to die
    sleep 2
    
    # Double check ffmpeg is dead
    if pgrep ffmpeg >/dev/null; then
        echo "Force killing ffmpeg..."
        pkill -9 ffmpeg
    fi
}

last_restart_modified=$(stat -c %Y /config/restart.txt)
last_shutdown_modified=$(stat -c %Y /config/shutdown.txt)

# Initial start
start_app

while true; do
    sleep 1
    new_restart_modified=$(stat -c %Y /config/restart.txt)
    new_shutdown_modified=$(stat -c %Y /config/shutdown.txt)
    
    if [ "$new_restart_modified" != "$last_restart_modified" ]; then
        echo "Restart triggered..."
        kill_app
        start_app
        last_restart_modified=$new_restart_modified
    elif [ "$new_shutdown_modified" != "$last_shutdown_modified" ]; then
        echo "Shutdown triggered..."
        kill_app
        exit 0
    fi
done