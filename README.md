# flux-fleet-scanner
**One‑line description:** Python tool that scans, validates, and reports on the Cocapn Fleet components.

## What it does
`flux-fleet-scanner` walks through the various Flux sub‑projects (e.g., `flux-a2a-prototype`, `flux-adaptive-opcodes`, `flux-bytecode-diff`, …) and:

- Detects missing or malformed configuration files (`.env`, CI scripts, etc.).
- Checks version compatibility across the fleet.
- Generates a concise health report (warnings, errors, and suggestions).
- Outputs results in JSON and human‑readable tables for CI pipelines.

## Installation
```bash
# Clone the repo
git clone https://github.com/SuperInstanceOrg/flux-fleet-scanner.git
cd flux-fleet-scanner

# Create a clean Python environment
python -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

## Configuration
Copy the example environment file and edit as needed:

```bash
cp .env.example .env
# Edit .env to set any required API keys, paths, etc.
```

## Usage
```bash
# Run the scanner on the entire repository
python -m flux_fleet_scanner

# Scan a specific sub‑project
python -m flux_fleet_scanner --path flux-coop-runtime

# Export results
python -m flux_fleet_scanner --output report.json
```

For advanced options, see `DOCKSIDE-EXAM.md`.

## Contributing
- Fork the repo and create a feature branch.  
- Follow the existing code style (PEP 8).  
- Submit a pull request with a clear description and tests.

## License
Distributed under the terms of the **MIT License** (see `LICENSE`).