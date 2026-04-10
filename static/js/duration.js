(function() {
    window.DurationTracker = DurationTracker;

    function DurationTracker(code) {
        this.code = code;
        this.startTime = Date.now();
        this.totalDuration = 0;
        this.isPaused = false;
        this.pauseStart = null;
        this.heartbeatInterval = null;

        this.bindEvents();
        this.startHeartbeat();
    }

    DurationTracker.prototype.bindEvents = function() {
        var self = this;

        document.addEventListener('visibilitychange', function() {
            self.handleVisibilityChange();
        });

        window.addEventListener('beforeunload', function() {
            self.reportDuration();
        });

        window.addEventListener('pagehide', function() {
            self.reportDuration();
        });
    };

    DurationTracker.prototype.handleVisibilityChange = function() {
        if (document.hidden) {
            this.isPaused = true;
            this.pauseStart = Date.now();
        } else {
            if (this.pauseStart) {
                this.totalDuration += Date.now() - this.pauseStart;
                this.pauseStart = null;
            }
            this.isPaused = false;
        }
    };

    DurationTracker.prototype.startHeartbeat = function() {
        var self = this;
        this.heartbeatInterval = setInterval(function() {
            self.sendHeartbeat();
        }, 30000);
    };

    DurationTracker.prototype.reportDuration = function() {
        var finalDuration = this.totalDuration;
        if (!this.isPaused && this.pauseStart === null) {
            finalDuration += Date.now() - this.startTime;
        }

        this.sendDuration(Math.floor(finalDuration / 1000));
    };

    DurationTracker.prototype.sendHeartbeat = function() {
        var self = this;
        if (!this.isPaused) {
            fetch('/api/visits/heartbeat', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    code: this.code,
                    timestamp: Date.now()
                })
            }).catch(function(err) {
                console.error('Heartbeat failed:', err);
            });
        }
    };

    DurationTracker.prototype.sendDuration = function(seconds) {
        fetch('/api/visits/duration', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                code: this.code,
                duration: seconds
            })
        }).catch(function(err) {
            console.error('Duration report failed:', err);
        });
    };

    DurationTracker.prototype.destroy = function() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
        }
        this.reportDuration();
    };
})();
