(function() {
    let currentShortlinksPage = 1;
    let currentVisitsPage = 1;
    let map = null;

    const API_BASE = '';

    document.addEventListener('DOMContentLoaded', function() {
        initNavigation();
        initModal();
        loadStats();
        loadShortlinks();
    });

    function initNavigation() {
        document.querySelectorAll('.nav-item').forEach(item => {
            item.addEventListener('click', function(e) {
                e.preventDefault();
                const page = this.dataset.page;
                switchPage(page);
            });
        });
    }

    function switchPage(page) {
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.toggle('active', item.dataset.page === page);
        });
        document.querySelectorAll('.page-section').forEach(section => {
            section.classList.toggle('active', section.id === page + 'Page');
        });

        const titles = { shortlinks: '短链接管理', visits: '访问记录', stats: '数据统计' };
        document.getElementById('pageTitle').textContent = titles[page] || '管理后台';

        if (page === 'shortlinks') {
            loadShortlinks();
        } else if (page === 'visits') {
            loadVisitCodeFilter();
            loadVisits();
        } else if (page === 'stats') {
            loadStats();
        }
    }

    function initModal() {
        const modal = document.getElementById('createModal');
        document.getElementById('createBtn').addEventListener('click', () => {
            modal.classList.add('active');
        });
        document.getElementById('closeModal').addEventListener('click', () => {
            modal.classList.remove('active');
        });
        document.getElementById('cancelBtn').addEventListener('click', () => {
            modal.classList.remove('active');
        });
        document.getElementById('submitBtn').addEventListener('click', createShortlink);
        
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.classList.remove('active');
            }
        });
    }

    async function createShortlink() {
        const url = document.getElementById('originalUrl').value;
        if (!url) {
            alert('请输入URL地址');
            return;
        }

        try {
            const response = await fetch(API_BASE + '/api/shortlinks/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url })
            });
            const result = await response.json();
            if (result.code === 0) {
                document.getElementById('createModal').classList.remove('active');
                document.getElementById('originalUrl').value = '';
                alert('创建成功！短链接: /' + result.data.code);
                loadShortlinks();
                loadStats();
            } else {
                alert('错误: ' + result.message);
            }
        } catch (error) {
            alert('创建失败');
        }
    }

    async function loadStats() {
        try {
            const response = await fetch(API_BASE + '/api/admin/stats');
            const result = await response.json();
            if (result.code === 0) {
                const data = result.data;
                document.getElementById('statVisits').textContent = data.total_visits || 0;
                document.getElementById('statDuration').textContent = (data.total_duration || 0) + 's';
                
                const linksResponse = await fetch(API_BASE + '/api/admin/shortlinks?page_size=1');
                const linksResult = await linksResponse.json();
                document.getElementById('statLinks').textContent = linksResult.data?.total || 0;

                renderVisitTrend(data.visit_trend || []);
                renderDeviceDistribution(data.device_distribution || {});
                renderGeoMap(data.geo_distribution || []);
            }
        } catch (error) {
            console.error('Failed to load stats:', error);
        }
    }

    async function loadShortlinks(page = 1) {
        currentShortlinksPage = page;
        try {
            const response = await fetch(API_BASE + '/api/admin/shortlinks?page=' + page + '&page_size=20');
            const result = await response.json();
            if (result.code === 0) {
                renderShortlinksTable(result.data.items || []);
                renderShortlinksPagination(result.data.total, result.data.page_size);
            }
        } catch (error) {
            console.error('Failed to load shortlinks:', error);
        }
    }

    function renderShortlinksTable(shortlinks) {
        const tbody = document.getElementById('shortlinksTableBody');
        if (!shortlinks || shortlinks.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="no-data">暂无数据</td></tr>';
            return;
        }

        tbody.innerHTML = shortlinks.map(s => {
            const status = s.is_disabled ? '<span class="badge warning">已禁用</span>' :
                          s.is_deleted ? '<span class="badge danger">已删除</span>' :
                          '<span class="badge success">正常</span>';
            const shortURL = window.location.origin + '/' + s.code;
            return `
                <tr>
                    <td><a href="/${s.code}" target="_blank">/${s.code}</a></td>
                    <td class="url-cell" title="${s.original_url}">${truncateUrl(s.original_url)}</td>
                    <td>${s.total_visits || 0}</td>
                    <td>${formatDate(s.created_at)}</td>
                    <td>${status}</td>
                    <td>
                        <button class="btn btn-secondary btn-small" onclick="copyLink('${shortURL}')">复制</button>
                        <button class="btn btn-secondary btn-small" onclick="window.open('/admin/visits?code=${s.code}')">记录</button>
                        <button class="btn btn-danger btn-small" onclick="deleteShortlink(${s.id})">删除</button>
                    </td>
                </tr>
            `;
        }).join('');
    }

    function renderShortlinksPagination(total, pageSize) {
        const totalPages = Math.ceil(total / pageSize);
        const pagination = document.getElementById('shortlinksPagination');
        if (totalPages <= 1) {
            pagination.innerHTML = '';
            return;
        }

        let html = '';
        for (let i = 1; i <= totalPages; i++) {
            html += `<button class="${i === currentShortlinksPage ? 'active' : ''}" onclick="loadShortlinks(${i})">${i}</button>`;
        }
        pagination.innerHTML = html;
    }

    async function loadVisits(page = 1) {
        currentVisitsPage = page;
        const code = document.getElementById('visitCodeFilter').value;
        try {
            const url = API_BASE + '/api/visits?page=' + page + '&page_size=50' + (code ? '&code=' + code : '');
            const response = await fetch(url);
            const result = await response.json();
            if (result.code === 0) {
                renderVisitsTable(result.data.items || []);
                renderVisitsPagination(result.data.total, result.data.page_size);
            }
        } catch (error) {
            console.error('Failed to load visits:', error);
        }
    }

    function renderVisitsTable(visits) {
        const tbody = document.getElementById('visitsTableBody');
        if (!visits || visits.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="no-data">暂无访问记录</td></tr>';
            return;
        }

        tbody.innerHTML = visits.map(v => {
            const location = [v.country, v.city].filter(Boolean).join(' ') || '-';
            const coords = v.latitude && v.longitude ?
                `<a href="https://maps.google.com/?q=${v.latitude},${v.longitude}" target="_blank">${v.latitude.toFixed(4)}, ${v.longitude.toFixed(4)}</a>` : '-';
            const device = [v.os_type, v.browser_type].filter(Boolean).join(' / ') || '-';
            const geoStatus = v.geo_status === 'granted' ? '<span class="badge success">GPS</span>' :
                             v.geo_status === 'denied' ? '<span class="badge warning">已拒绝</span>' :
                             '<span class="badge danger">失败</span>';
            return `
                <tr>
                    <td>${v.ip_address}</td>
                    <td>${location} ${geoStatus}</td>
                    <td>${coords}</td>
                    <td>${device}</td>
                    <td>${v.visit_duration || 0}s</td>
                    <td>${formatDate(v.visit_time)}</td>
                </tr>
            `;
        }).join('');
    }

    function renderVisitsPagination(total, pageSize) {
        const totalPages = Math.ceil(total / pageSize);
        const pagination = document.getElementById('visitsPagination');
        if (totalPages <= 1) {
            pagination.innerHTML = '';
            return;
        }

        let html = '';
        for (let i = 1; i <= totalPages; i++) {
            html += `<button class="${i === currentVisitsPage ? 'active' : ''}" onclick="loadVisits(${i})">${i}</button>`;
        }
        pagination.innerHTML = html;
    }

    async function loadVisitCodeFilter() {
        try {
            const response = await fetch(API_BASE + '/api/admin/shortlinks?page_size=100');
            const result = await response.json();
            if (result.code === 0) {
                const select = document.getElementById('visitCodeFilter');
                select.innerHTML = '<option value="">全部链接</option>';
                (result.data.items || []).forEach(s => {
                    select.innerHTML += `<option value="${s.code}">/${s.code}</option>`;
                });
                select.addEventListener('change', () => loadVisits());
            }
        } catch (error) {
            console.error('Failed to load shortlinks for filter:', error);
        }
    }

    async function deleteShortlink(id) {
        if (!confirm('确定要删除这个短链接吗？')) return;
        try {
            const response = await fetch(API_BASE + '/api/admin/shortlinks/' + id, { method: 'DELETE' });
            const result = await response.json();
            if (result.code === 0) {
                loadShortlinks(currentShortlinksPage);
                loadStats();
            } else {
                alert('错误: ' + result.message);
            }
        } catch (error) {
            alert('删除失败');
        }
    }

    function copyLink(url) {
        navigator.clipboard.writeText(url).then(() => {
            alert('链接已复制到剪贴板');
        }).catch(() => {
            prompt('请复制链接:', url);
        });
    }

    function renderVisitTrend(trend) {
        const container = document.getElementById('visitTrendChart');
        if (!trend || trend.length === 0) {
            container.innerHTML = '<p class="no-data">暂无数据</p>';
            return;
        }

        const maxCount = Math.max(...trend.map(t => t.count), 1);
        container.innerHTML = '<div class="bar-chart-container">' + trend.map(item => {
            const height = Math.max((item.count / maxCount) * 150, 4);
            return `
                <div class="bar-item">
                    <div class="bar" style="height: ${height}px">
                        <span class="bar-value">${item.count}</span>
                    </div>
                    <div class="bar-label">${item.date.slice(5)}</div>
                </div>
            `;
        }).join('') + '</div>';
    }

    function renderDeviceDistribution(dist) {
        const container = document.getElementById('deviceDistChart');
        const total = Object.values(dist).reduce((a, b) => a + b, 0);
        if (total === 0) {
            container.innerHTML = '<p class="no-data">暂无数据</p>';
            return;
        }

        const colors = { pc: '#667eea', mobile: '#764ba2', tablet: '#f093fb' };
        const labels = { pc: '电脑', mobile: '手机', tablet: '平板' };
        container.innerHTML = '<div class="pie-chart-container">' + Object.entries(dist).map(([device, count]) => {
            const percent = ((count / total) * 100).toFixed(1);
            return `
                <div class="pie-item">
                    <div class="pie-color" style="background: ${colors[device] || '#999'}"></div>
                    <div class="pie-label">${labels[device] || device}: ${count} (${percent}%)</div>
                </div>
            `;
        }).join('') + '</div>';
    }

    function renderGeoMap(geoPoints) {
        if (!map) {
            map = L.map('map').setView([35, 105], 4);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '© OpenStreetMap'
            }).addTo(map);
        }

        map.eachLayer(layer => {
            if (layer instanceof L.Circle || layer instanceof L.Marker) {
                map.removeLayer(layer);
            }
        });

        if (!geoPoints || geoPoints.length === 0) {
            return;
        }

        geoPoints.forEach(point => {
            if (point.lat && point.lng) {
                L.circle([point.lat, point.lng], {
                    radius: Math.max(point.count * 200, 500),
                    color: '#667eea',
                    fillColor: '#667eea',
                    fillOpacity: 0.5
                }).addTo(map).bindPopup(`访问次数: ${point.count}`);
            }
        });

        if (geoPoints.length === 1) {
            map.setView([geoPoints[0].lat, geoPoints[0].lng], 12);
        } else {
            const bounds = L.latLngBounds(geoPoints.map(p => [p.lat, p.lng]));
            map.fitBounds(bounds, { padding: [50, 50] });
        }
    }

    document.getElementById('exportBtn').addEventListener('click', function() {
        const code = document.getElementById('visitCodeFilter').value;
        window.open(API_BASE + '/api/admin/export?code=' + code, '_blank');
    });

    function truncateUrl(url, maxLen = 40) {
        if (!url) return '-';
        if (url.length <= maxLen) return url;
        return url.substring(0, maxLen) + '...';
    }

    function formatDate(dateStr) {
        if (!dateStr) return '-';
        const date = new Date(dateStr);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    window.loadShortlinks = loadShortlinks;
    window.loadVisits = loadVisits;
    window.deleteShortlink = deleteShortlink;
    window.copyLink = copyLink;
})();
