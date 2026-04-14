# flux-fleet-scanner
**Python tool for scanning and analyzing Flux fleet data** – a component of the **Cocapn Fleet** (github.com/SuperInstance).

## Description
`flux-fleet-scanner` walks through fleet repositories, extracts metadata, validates bytecode compatibility, and generates concise reports. It integrates with the other Flux modules (e.g., `flux-cooperative-intelligence`, `flux-evolution`) to provide a holistic view of fleet health.

## Usage
```bash
# Clone the repo
git clone https://github.com/SuperInstance/flux-fleet-scanner.git
cd flux-fleet-scanner

# Install dependencies
pip install -r requirements.txt

# Run the scanner (default scans the `download/` directory)
python -m flux_fleet_scanner --path download/ --output report.json
```
*Optional flags:* `--verbose`, `--filter <module>`, `--format yaml|json`.

## Related Projects
- **Cocapn Fleet** – overall fleet orchestration: https://github.com/SuperInstance  
- **flux-cooperative-intelligence** – AI‑driven fleet coordination: https://github.com/SuperInstance/flux-cooperative-intelligence  
- **flux-evolution** – version‑migration utilities: https://github.com/SuperInstance/flux-evolution  

## License
Distributed under the terms of the [MIT License](LICENSE).