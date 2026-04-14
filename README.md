# flux-fleet-scanner
**Python scanner for the Cocapn Fleet**  

Part of the Cocapn Fleet ecosystem (github.com/SuperInstance).

## Description
`flux-fleet-scanner` inspects, validates, and reports on fleet components such as flux‑bytecode, adaptive opcodes, and cooperative intelligence modules. It helps maintain consistency across the fleet’s diverse repositories.

## Usage
```bash
# Clone the repo
git clone https://github.com/SuperInstance/flux-fleet-scanner.git
cd flux-fleet-scanner

# Install dependencies
pip install -r requirements.txt

# Run the scanner (default scans the `download/` directory)
python -m flux_fleet_scanner --path ./download
```
Optional flags:
- `--config .env` – load environment variables.
- `--output report.json` – write results to a JSON file.

## Related Projects
- [flux-cooperative-intelligence](https://github.com/SuperInstance/flux-cooperative-intelligence) – AI‑driven fleet coordination.
- [flux-a2a-prototype](https://github.com/SuperInstance/flux-a2a-prototype) – A2A communication layer.
- [flux-conformance](https://github.com/SuperInstance/flux-conformance) – Test suite for fleet standards.

## Contributing
Contributions are welcome. Please open issues or submit pull requests against the `main` branch.

## License
Distributed under the terms of the [LICENSE](LICENSE).