(function() {
    window.Tracker = {
        code: null,
        destinationUrl: null,
        geoSubmitted: false,
        durationTracker: null,

        init: function(code, destUrl) {
            this.code = code;
            this.destinationUrl = destUrl;
            this.startGeoTracking();
            this.startDurationTracking();
        },

        startGeoTracking: function() {
            if (navigator.geolocation) {
                navigator.geolocation.getCurrentPosition(
                    this.onGeoSuccess.bind(this),
                    this.onGeoError.bind(this),
                    {
                        enableHighAccuracy: true,
                        timeout: 8000,
                        maximumAge: 300000
                    }
                );
            } else {
                this.submitGeo(null, null, 'gps', 'unsupported');
            }
        },

        onGeoSuccess: function(position) {
            this.submitGeo(
                position.coords.latitude,
                position.coords.longitude,
                'gps',
                'granted'
            );
        },

        onGeoError: function(error) {
            var status = 'failed';
            switch(error.code) {
                case error.PERMISSION_DENIED:
                    status = 'denied';
                    break;
                case error.POSITION_UNAVAILABLE:
                    status = 'unavailable';
                    break;
                case error.TIMEOUT:
                    status = 'timeout';
                    break;
            }
            this.submitGeo(null, null, 'gps', status);
        },

        submitGeo: function(lat, lng, precision, status) {
            if (this.geoSubmitted) return;
            this.geoSubmitted = true;

            var data = {
                code: this.code,
                latitude: lat,
                longitude: lng,
                geo_precision: precision,
                geo_status: status,
                visit_duration: 0,
                page_loaded: true
            };

            if (navigator.sendBeacon) {
                navigator.sendBeacon('/api/visits/submit', JSON.stringify(data));
            } else {
                fetch('/api/visits/submit', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(data)
                }).catch(function(err) {
                    console.error('Failed to submit geo data:', err);
                });
            }
        },

        startDurationTracking: function() {
            this.durationTracker = new DurationTracker(this.code);
        }
    };

    function DurationTracker(code) {
        this.code = code;
        this.startTime = Date.now();
        this.totalDuration = 0;
        this.isPaused = false;
        this.pauseStart = null;

        var self = this;

        document.addEventListener('visibilitychange', function() {
            self.handleVisibilityChange();
        });

        window.addEventListener('beforeunload', function() {
            self.reportDuration();
        });

        this.heartbeatInterval = setInterval(function() {
            self.sendHeartbeat();
        }, 30000);
    }

    DurationTracker.prototype.handleVisibilityChange = function() {
        var self = this;
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

    DurationTracker.prototype.reportDuration = function() {
        var finalDuration = this.totalDuration;
        if (!this.isPaused && this.pauseStart === null) {
            finalDuration += Date.now() - this.startTime;
        }

        var data = {
            code: this.code,
            duration: Math.floor(finalDuration / 1000)
        };

        if (navigator.sendBeacon) {
            navigator.sendBeacon('/api/visits/duration', JSON.stringify(data));
        } else {
            fetch('/api/visits/duration', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(data)
            });
        }
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
})();
