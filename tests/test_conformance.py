"""
Tests for the FLUX Conformance Test Suite.

Tests the test vectors and runner function from conformance_tests.py.
"""

import sys
import os
import pytest

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "download"))
from conformance_tests import TEST_VECTORS, run_conformance_tests


# ===========================================================================
# Test Vector Structure Tests
# ===========================================================================

class TestConformanceTestVectors:
    """Verify the test vector data structure is correct."""

    def test_vectors_exist(self):
        assert len(TEST_VECTORS) > 0

    def test_all_vectors_have_name(self):
        for v in TEST_VECTORS:
            assert "name" in v
            assert isinstance(v["name"], str)
            assert len(v["name"]) > 0

    def test_all_vectors_have_category(self):
        for v in TEST_VECTORS:
            assert "category" in v
            assert isinstance(v["category"], str)

    def test_all_vectors_have_expected(self):
        for v in TEST_VECTORS:
            assert "expected" in v

    def test_categories_are_valid(self):
        valid_categories = {"control", "data", "arithmetic", "comparison",
                           "stack", "logic", "complex"}
        for v in TEST_VECTORS:
            assert v["category"] in valid_categories

    def test_no_crash_vectors_exist(self):
        """There should be at least one 'no_crash' test."""
        no_crash = [v for v in TEST_VECTORS if v["expected"] == "no_crash"]
        assert len(no_crash) >= 1

    def test_register_check_vectors_exist(self):
        """There should be tests that check specific register values."""
        reg_tests = [v for v in TEST_VECTORS
                     if isinstance(v["expected"], dict) and "register" in v["expected"]]
        assert len(reg_tests) >= 5

    def test_arithmetic_tests_exist(self):
        arith = [v for v in TEST_VECTORS if v["category"] == "arithmetic"]
        assert len(arith) >= 3

    def test_logic_tests_exist(self):
        logic = [v for v in TEST_VECTORS if v["category"] == "logic"]
        assert len(logic) >= 2

    def test_comparison_tests_exist(self):
        comp = [v for v in TEST_VECTORS if v["category"] == "comparison"]
        assert len(comp) >= 1

    def test_stack_tests_exist(self):
        stack = [v for v in TEST_VECTORS if v["category"] == "stack"]
        assert len(stack) >= 1

    def test_data_tests_exist(self):
        data = [v for v in TEST_VECTORS if v["category"] == "data"]
        assert len(data) >= 1

    def test_complex_tests_exist(self):
        complex_tests = [v for v in TEST_VECTORS if v["category"] == "complex"]
        assert len(complex_tests) >= 1

    def test_bytecode_vectors_have_bytecode(self):
        """Non-complex vectors should have bytecode."""
        for v in TEST_VECTORS:
            if v["category"] != "complex" and v["bytecode"] is None:
                assert False, f"Non-complex test '{v['name']}' has no bytecode"

    def test_notes_exist(self):
        """All vectors should have notes explaining the test."""
        for v in TEST_VECTORS:
            assert "notes" in v
            assert isinstance(v["notes"], str)

    def test_bytecode_values_are_ints(self):
        for v in TEST_VECTORS:
            if v["bytecode"] is not None:
                for byte in v["bytecode"]:
                    assert isinstance(byte, int)


# ===========================================================================
# Runner Function Tests
# ===========================================================================

class TestConformanceRunner:
    """Test the run_conformance_tests runner function."""

    def _simple_runner(self, bytecode):
        """A minimal mock VM that returns all registers as 0."""
        return {"registers": {}, "crashed": False}

    def test_runner_returns_results_dict(self):
        results = run_conformance_tests(self._simple_runner)
        assert "passed" in results
        assert "failed" in results
        assert "results" in results

    def test_runner_all_no_crash_pass(self):
        """no_crash tests should pass with non-crashing runner."""
        results = run_conformance_tests(self._simple_runner)
        no_crash_results = [r for r in results["results"] if "no_crash" in r["name"].lower() or "NOP" in r["name"].upper() or "HALT" in r["name"].upper()]
        # At minimum NOP and HALT should pass
        assert results["passed"] >= 2

    def test_runner_skips_none_bytecode(self):
        """Tests with None bytecode should be skipped."""
        results = run_conformance_tests(self._simple_runner)
        skipped = [r for r in results["results"] if r.get("status") == "SKIPPED"]
        assert len(skipped) >= 1

    def test_runner_counts_pass_fail(self):
        results = run_conformance_tests(self._simple_runner)
        assert results["passed"] + results["failed"] + len(
            [r for r in results["results"] if r.get("status") == "SKIPPED"]
        ) == len(TEST_VECTORS)

    def test_runner_handles_crash(self):
        """Runner should detect crashes."""
        def crashing_runner(bytecode):
            return {"registers": {}, "crashed": True}

        results = run_conformance_tests(crashing_runner)
        no_crash_tests = [v for v in TEST_VECTORS if v["expected"] == "no_crash"]
        assert results["failed"] >= len(no_crash_tests)

    def test_runner_handles_register_mismatch(self):
        """Runner should detect register value mismatches."""
        def wrong_runner(bytecode):
            return {"registers": {0: 999, 1: 999, 2: 999}, "crashed": False}

        results = run_conformance_tests(wrong_runner)
        assert results["failed"] > 0

    def test_runner_handles_exception(self):
        """Runner should handle exceptions from the VM."""
        def error_runner(bytecode):
            raise RuntimeError("VM exploded")

        results = run_conformance_tests(error_runner)
        assert results["failed"] > 0
        error_results = [r for r in results["results"] if r.get("status") == "ERROR"]
        assert len(error_results) > 0

    def test_runner_results_have_name(self):
        results = run_conformance_tests(self._simple_runner)
        for r in results["results"]:
            assert "name" in r

    def test_runner_results_have_status(self):
        results = run_conformance_tests(self._simple_runner)
        for r in results["results"]:
            assert "status" in r
            assert r["status"] in ["PASS", "FAIL", "SKIPPED", "ERROR"]


# ===========================================================================
# Specific Test Vector Content Tests
# ===========================================================================

class TestSpecificVectors:
    """Verify specific test vector content."""

    def test_nop_halt_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("NOP" in n for n in names)

    def test_halt_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("HALT" in n for n in names)

    def test_add_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("ADD" in n for n in names)

    def test_sub_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("SUB" in n for n in names)

    def test_mul_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("MUL" in n for n in names)

    def test_mod_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("MOD" in n for n in names)

    def test_and_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("AND" in n for n in names)

    def test_or_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("OR" in n.upper() or "Or" in n for n in names)

    def test_xor_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("XOR" in n.upper() or "Xor" in n for n in names)

    def test_gcd_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("GCD" in n or "gcd" in n for n in names)

    def test_fibonacci_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("Fibonacci" in n or "fibonacci" in n for n in names)

    def test_register_overlap_tests_exist(self):
        overlap_tests = [v for v in TEST_VECTORS if "overlap" in v.get("notes", "").lower()]
        assert len(overlap_tests) >= 2

    def test_cmp_eq_equal_value_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("equal" in n.lower() and "CMP" in n for n in names)

    def test_cmp_eq_unequal_value_exists(self):
        names = [v["name"] for v in TEST_VECTORS]
        assert any("unequal" in n.lower() and "CMP" in n for n in names)

    def test_value_neq_zero_flag_exists(self):
        """Verify value_neq_zero expected key is used."""
        neq_zero = [v for v in TEST_VECTORS
                     if isinstance(v["expected"], dict) and "value_neq_zero" in v["expected"]]
        assert len(neq_zero) >= 1
