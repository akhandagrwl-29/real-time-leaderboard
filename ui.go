package main

import "html/template"

var mainUITmpl = template.Must(template.New("mainui").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Leaderboard System</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body { padding: 2em; }
        .tab-content { margin-top: 2em; }
        pre { background: #f8f9fa; padding: 1em; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="mb-4">Leaderboard System</h1>
        <ul class="nav nav-tabs" id="uiTab" role="tablist">
            <li class="nav-item" role="presentation">
                <button class="nav-link active" id="submit-tab" data-bs-toggle="tab" data-bs-target="#submit" type="button" role="tab">Submit Score</button>
            </li>
            <li class="nav-item" role="presentation">
                <button class="nav-link" id="leaderboard-tab" data-bs-toggle="tab" data-bs-target="#leaderboard" type="button" role="tab">Leaderboard</button>
            </li>
            <li class="nav-item" role="presentation">
                <button class="nav-link" id="rank-tab" data-bs-toggle="tab" data-bs-target="#rank" type="button" role="tab">User Rank</button>
            </li>
        </ul>
        <div class="tab-content">
            <!-- Submit Score Tab -->
            <div class="tab-pane fade show active" id="submit" role="tabpanel">
                <form id="submitForm" class="mt-4">
                    <div class="mb-3">
                        <label>User</label>
                        <input name="user" class="form-control" required>
                    </div>
                    <div class="mb-3">
                        <label>Game</label>
                        <input name="game" class="form-control" required>
                    </div>
                    <div class="mb-3">
                        <label>Score</label>
                        <input name="score" type="number" step="any" class="form-control" required>
                    </div>
                    <button type="submit" class="btn btn-primary">Submit</button>
                </form>
                <pre id="submitResult" class="mt-3"></pre>
            </div>
            <!-- Leaderboard Tab -->
            <div class="tab-pane fade" id="leaderboard" role="tabpanel">
                <form id="lbForm" class="mt-4 row g-3 align-items-end">
                    <div class="col-auto">
                        <label>Game</label>
                        <input name="game" class="form-control" required>
                    </div>
                    <div class="col-auto">
                        <label>Top N</label>
                        <input name="n" type="number" value="10" min="1" class="form-control">
                    </div>
                    <div class="col-auto">
                        <button type="submit" class="btn btn-success">Show</button>
                    </div>
                </form>
                <table class="table table-striped mt-3" id="lbTable" style="display:none;">
                    <thead><tr><th>Rank</th><th>User</th><th>Score</th></tr></thead>
                    <tbody></tbody>
                </table>
                <pre id="lbResult"></pre>
            </div>
            <!-- User Rank Tab -->
            <div class="tab-pane fade" id="rank" role="tabpanel">
                <form id="rankForm" class="mt-4 row g-3 align-items-end">
                    <div class="col-auto">
                        <label>User</label>
                        <input name="user" class="form-control" required>
                    </div>
                    <div class="col-auto">
                        <label>Game</label>
                        <input name="game" class="form-control" required>
                    </div>
                    <div class="col-auto">
                        <button type="submit" class="btn btn-info">Show</button>
                    </div>
                </form>
                <pre id="rankResult" class="mt-3"></pre>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js"></script>
    <script>
    // Submit Score
    document.getElementById('submitForm').onsubmit = async function(e) {
        e.preventDefault();
        const user = this.user.value;
        const game = this.game.value;
        const score = parseFloat(this.score.value);
        const res = await fetch('/submit', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({user, game, score})
        });
        document.getElementById('submitResult').textContent = await res.text();
    };
    // Leaderboard
    document.getElementById('lbForm').onsubmit = async function(e) {
        e.preventDefault();
        const game = this.game.value;
        const n = this.n.value;
        const res = await fetch('/leaderboard?game=' + encodeURIComponent(game) + '&n=' + n);
        if (res.ok) {
        const data = await res.json();
        const tbody = document.querySelector('#lbTable tbody');
        tbody.innerHTML = '';

        data.forEach(row => {
            tbody.innerHTML +=
                '<tr>' +
                '<td>' + row.rank + '</td>' +
                '<td>' + row.user + '</td>' +
                '<td>' + row.score + '</td>' +
                '</tr>';
        });

            document.getElementById('lbTable').style.display = '';
            document.getElementById('lbResult').textContent = '';
        } else {
            document.getElementById('lbTable').style.display = 'none';
            document.getElementById('lbResult').textContent = await res.text();
        }
    };
    // User Rank
    document.getElementById('rankForm').onsubmit = async function(e) {
        e.preventDefault();
        const user = this.user.value;
        const game = this.game.value;
        const res = await fetch('/rank?user=' + encodeURIComponent(user) + '&game=' + encodeURIComponent(game));
        if (res.ok) {
            document.getElementById('rankResult').textContent = JSON.stringify(await res.json(), null, 2);
        } else {
            document.getElementById('rankResult').textContent = await res.text();
        }
    };
    </script>
</body>
</html>
`))
