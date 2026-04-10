(function() {
    window.GeoTracker = {
        code: null,
        submitted: false,

        init: function(code) {
            this.code = code;
            this.requestLocation();
        },

        requestLocation: function() {
            if (navigator.geolocation) {
                navigator.geolocation.getCurrentPosition(
                    this.onSuccess.bind(this),
                    this.onError.bind(this),
                    {
                        enableHighAccuracy: true,
                        timeout: 10000,
                        maximumAge: 300000
                    }
                );
            } else {
                this.submit(null, null, 'gps', 'unsupported');
            }
        },

        onSuccess: function(position) {
            this.submit(
                position.coords.latitude,
                position.coords.longitude,
                'gps',
                'granted'
            );
        },

        onError: function(error) {
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
            this.submit(null, null, 'gps', status);
        },

        submit: function(lat, lng, precision, geoStatus) {
            if (this.submitted) return;
            this.submitted = true;

            var data = {
                code: this.code,
                latitude: lat,
                longitude: lng,
                geo_precision: precision,
                geo_status: geoStatus,
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
                    console.error('Geo submission failed:', err);
                });
            }
        }
    };
})();
