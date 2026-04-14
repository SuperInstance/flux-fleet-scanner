"""
Tests for the fleet discovery module (lib/discovery.py).

Tests use mocking to avoid actual GitHub API calls.
"""

import sys
import os
import json
import pytest
import importlib
from unittest.mock import patch, MagicMock, call
import urllib.request as real_urllib

# Import the discovery module directly
_discovery_path = os.path.join(os.path.dirname(__file__), "..", "repos", "greenhorn-runtime", "lib", "discovery.py")
import importlib.util
_spec = importlib.util.spec_from_file_location("discovery", _discovery_path)
discovery = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(discovery)

discover_fleet_repos = discovery.discover_fleet_repos
find_vessel_repos = discovery.find_vessel_repos
find_task_repos = discovery.find_task_repos
GITHUB_TOKEN = discovery.GITHUB_TOKEN
API = discovery.API


class TestDiscoveryConfig:
    """Test module configuration."""

    def test_api_url(self):
        assert API == "https://api.github.com"

    def test_github_token_is_string(self):
        assert isinstance(GITHUB_TOKEN, str)


class TestDiscoverFleetRepos:
    """Test discover_fleet_repos with mocked API."""

    def _mock_repos(self, count=5, include_forks=False):
        repos = []
        for i in range(count):
            repo = {
                "name": f"repo-{i}",
                "full_name": f"SuperInstance/repo-{i}",
                "fork": i == 3 if include_forks else False,
                "updated_at": "2026-04-12T00:00:00Z",
            }
            repos.append(repo)
        return repos

    def _make_mock_urlopen(self, data):
        mock_response = MagicMock()
        mock_response.read.return_value = json.dumps(data).encode()
        mock_response.__enter__ = MagicMock(return_value=mock_response)
        mock_response.__exit__ = MagicMock(return_value=False)
        return mock_response

    def test_returns_non_fork_repos(self):
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen(self._mock_repos(5))
            repos = discover_fleet_repos()
            assert len(repos) == 5
            assert all(not r.get("fork") for r in repos)

    def test_filters_out_forks(self):
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen(
                self._mock_repos(5, include_forks=True)
            )
            result = discover_fleet_repos()
            assert len(result) == 4

    def test_empty_response(self):
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen([])
            repos = discover_fleet_repos()
            assert repos == []

    def test_custom_owner(self):
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen([])
            mock_urllib.request.Request = real_urllib.Request
            discover_fleet_repos(owner="CustomOrg")
            req_obj = mock_urllib.request.urlopen.call_args[0][0]
            assert "CustomOrg" in req_obj.full_url

    def test_custom_per_page(self):
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen([])
            mock_urllib.request.Request = real_urllib.Request
            discover_fleet_repos(per_page=50)
            req_obj = mock_urllib.request.urlopen.call_args[0][0]
            assert "per_page=50" in req_obj.full_url

    def test_handles_api_error(self):
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.side_effect = Exception("API Error")
            with pytest.raises(Exception):
                discover_fleet_repos()

    def test_default_owner_is_superinstance(self):
        """Default owner should be SuperInstance."""
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen([])
            mock_urllib.request.Request = real_urllib.Request
            discover_fleet_repos()
            req_obj = mock_urllib.request.urlopen.call_args[0][0]
            assert "SuperInstance" in req_obj.full_url

    def test_default_per_page_is_100(self):
        """Default per_page should be 100."""
        with patch.object(discovery, "urllib") as mock_urllib:
            mock_urllib.request.urlopen.return_value = self._make_mock_urlopen([])
            mock_urllib.request.Request = real_urllib.Request
            discover_fleet_repos()
            req_obj = mock_urllib.request.urlopen.call_args[0][0]
            assert "per_page=100" in req_obj.full_url


class TestFindVesselRepos:
    """Test find_vessel_repos."""

    def test_filters_vessel_repos(self):
        with patch.object(discovery, "discover_fleet_repos") as mock_discover:
            mock_discover.return_value = [
                {"name": "superz-vessel"},
                {"name": "flux-runtime"},
                {"name": "oracle1-vessel"},
            ]
            vessels = find_vessel_repos()
            assert len(vessels) == 2
            assert all(r["name"].endswith("-vessel") for r in vessels)

    def test_no_vessel_repos(self):
        with patch.object(discovery, "discover_fleet_repos") as mock_discover:
            mock_discover.return_value = [
                {"name": "flux-runtime"},
                {"name": "flux-skills"},
            ]
            vessels = find_vessel_repos()
            assert vessels == []

    def test_all_vessel_repos(self):
        with patch.object(discovery, "discover_fleet_repos") as mock_discover:
            mock_discover.return_value = [
                {"name": "a-vessel"},
                {"name": "b-vessel"},
            ]
            vessels = find_vessel_repos()
            assert len(vessels) == 2


class TestFindTaskRepos:
    """Test find_task_repos."""

    def test_finds_task_repos(self):
        with patch.object(discovery, "discover_fleet_repos") as mock_discover, \
             patch.object(discovery, "urllib") as mock_urllib:
            mock_discover.return_value = [
                {"name": "repo-with-tasks"},
                {"name": "repo-without-tasks"},
            ]
            mock_urllib.request.urlopen.side_effect = [
                self._make_mock_urlopen_success(b"content"),
                Exception("Not found"),
            ]
            tasks = find_task_repos()
            assert len(tasks) == 1
            assert tasks[0]["name"] == "repo-with-tasks"

    def test_no_task_repos(self):
        with patch.object(discovery, "discover_fleet_repos") as mock_discover, \
             patch.object(discovery, "urllib") as mock_urllib:
            mock_discover.return_value = [{"name": "repo-no-tasks"}]
            mock_urllib.request.urlopen.side_effect = Exception("Not found")
            tasks = find_task_repos()
            assert tasks == []

    def _make_mock_urlopen_success(self, data):
        mock_response = MagicMock()
        mock_response.read.return_value = data
        mock_response.__enter__ = MagicMock(return_value=mock_response)
        mock_response.__exit__ = MagicMock(return_value=False)
        return mock_response
