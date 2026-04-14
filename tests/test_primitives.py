"""
Comprehensive tests for FLUX A2A Protocol Primitives.

Tests cover:
1. All 6 primitive types (Branch, Fork, CoIterate, Discuss, Synthesize, Reflect)
2. All sub-types (BranchBody, ForkMutation, ForkInherit, Participant, SynthesisSource)
3. All enums (BranchStrategy, MergeStrategy, ForkOnComplete, ForkConflictMode,
   SharedStateMode, DiscussFormat, SynthesisMethod, ReflectTarget)
4. JSON round-trip (to_dict -> from_dict)
5. Schema versioning ($schema field)
6. Confidence clamping
7. Default values
8. Registry and parse_primitive dispatch
9. Edge cases (empty, missing fields, unknown ops, meta dict)
10. Integration tests across all primitives
"""

import json
import pytest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "download"))
from primitives import (
    BranchPrimitive, BranchBody,
    ForkPrimitive, ForkInherit, ForkMutation,
    CoIteratePrimitive,
    DiscussPrimitive, Participant,
    SynthesizePrimitive, SynthesisSource,
    ReflectPrimitive,
    parse_primitive, PRIMITIVE_REGISTRY,
    BranchStrategy, MergeStrategy, ForkOnComplete, ForkConflictMode,
    SharedStateMode, DiscussFormat, SynthesisMethod, ReflectTarget,
    _clamp, _now, _uuid,
)


def round_trip(prim):
    """Serialize to dict, deserialize, return result."""
    return prim.__class__.from_dict(prim.to_dict())


# ===========================================================================
# Enum Tests
# ===========================================================================

class TestEnums:
    """Test all enum values are correct."""

    def test_branch_strategy_values(self):
        assert BranchStrategy.PARALLEL.value == "parallel"
        assert BranchStrategy.SEQUENTIAL.value == "sequential"
        assert BranchStrategy.COMPETITIVE.value == "competitive"

    def test_merge_strategy_values(self):
        assert MergeStrategy.CONSENSUS.value == "consensus"
        assert MergeStrategy.VOTE.value == "vote"
        assert MergeStrategy.BEST.value == "best"
        assert MergeStrategy.ALL.value == "all"
        assert MergeStrategy.WEIGHTED_CONFIDENCE.value == "weighted_confidence"
        assert MergeStrategy.FIRST_COMPLETE.value == "first_complete"
        assert MergeStrategy.LAST_WRITER_WINS.value == "last_writer_wins"
        assert MergeStrategy.CUSTOM.value == "custom"

    def test_fork_on_complete_values(self):
        assert ForkOnComplete.COLLECT.value == "collect"
        assert ForkOnComplete.DISCARD.value == "discard"
        assert ForkOnComplete.SIGNAL.value == "signal"
        assert ForkOnComplete.MERGE.value == "merge"

    def test_fork_conflict_mode_values(self):
        assert ForkConflictMode.PARENT_WINS.value == "parent_wins"
        assert ForkConflictMode.CHILD_WINS.value == "child_wins"
        assert ForkConflictMode.NEGOTIATE.value == "negotiate"

    def test_shared_state_mode_values(self):
        assert SharedStateMode.CONFLICT.value == "conflict"
        assert SharedStateMode.MERGE.value == "merge"
        assert SharedStateMode.PARTITIONED.value == "partitioned"
        assert SharedStateMode.ISOLATED.value == "isolated"

    def test_discuss_format_values(self):
        assert DiscussFormat.DEBATE.value == "debate"
        assert DiscussFormat.BRAINSTORM.value == "brainstorm"
        assert DiscussFormat.REVIEW.value == "review"
        assert DiscussFormat.NEGOTIATE.value == "negotiate"
        assert DiscussFormat.PEER_REVIEW.value == "peer_review"

    def test_synthesis_method_values(self):
        assert SynthesisMethod.MAP_REDUCE.value == "map_reduce"
        assert SynthesisMethod.ENSEMBLE.value == "ensemble"
        assert SynthesisMethod.CHAIN.value == "chain"
        assert SynthesisMethod.VOTE.value == "vote"
        assert SynthesisMethod.WEIGHTED_MERGE.value == "weighted_merge"
        assert SynthesisMethod.BEST_EFFORT.value == "best_effort"

    def test_reflect_target_values(self):
        assert ReflectTarget.STRATEGY.value == "strategy"
        assert ReflectTarget.PROGRESS.value == "progress"
        assert ReflectTarget.UNCERTAINTY.value == "uncertainty"
        assert ReflectTarget.CONFIDENCE.value == "confidence"
        assert ReflectTarget.ALL.value == "all"

    def test_enum_count(self):
        """Verify we have exactly 8 enums."""
        enums = [BranchStrategy, MergeStrategy, ForkOnComplete,
                 ForkConflictMode, SharedStateMode, DiscussFormat,
                 SynthesisMethod, ReflectTarget]
        assert len(enums) == 8


# ===========================================================================
# Helper Function Tests
# ===========================================================================

class TestHelpers:
    """Test internal helper functions."""

    def test_clamp_normal(self):
        assert _clamp(0.5) == 0.5

    def test_clamp_high(self):
        assert _clamp(1.5) == 1.0

    def test_clamp_low(self):
        assert _clamp(-0.5) == 0.0

    def test_clamp_zero(self):
        assert _clamp(0.0) == 0.0

    def test_clamp_one(self):
        assert _clamp(1.0) == 1.0

    def test_clamp_custom_range(self):
        assert _clamp(0.5, 0.0, 1.0) == 0.5
        assert _clamp(-1.0, -5.0, 5.0) == -1.0
        assert _clamp(10.0, -5.0, 5.0) == 5.0

    def test_clamp_string_input(self):
        """_clamp should convert string to float."""
        assert _clamp("0.5") == 0.5

    def test_now_returns_string(self):
        assert isinstance(_now(), str)

    def test_now_has_tz(self):
        """_now should include timezone info (UTC)."""
        result = _now()
        assert "+" in result or "Z" in result or "00:00" in result

    def test_uuid_returns_string(self):
        assert isinstance(_uuid(), str)

    def test_uuid_unique(self):
        assert _uuid() != _uuid()


# ===========================================================================
# BranchBody Tests
# ===========================================================================

class TestBranchBody:
    def test_default_values(self):
        bb = BranchBody()
        assert bb.label == ""
        assert bb.weight == 1.0
        assert bb.body == []
        assert bb.confidence == 1.0
        assert bb.meta == {}

    def test_weight_clamped_high(self):
        bb = BranchBody(weight=2.0)
        assert bb.weight == 1.0

    def test_weight_clamped_low(self):
        bb = BranchBody(weight=-1.0)
        assert bb.weight == 0.0

    def test_confidence_clamped_high(self):
        bb = BranchBody(confidence=1.5)
        assert bb.confidence == 1.0

    def test_confidence_clamped_low(self):
        bb = BranchBody(confidence=-0.5)
        assert bb.confidence == 0.0

    def test_to_dict_minimal(self):
        bb = BranchBody(label="test")
        d = bb.to_dict()
        assert d["label"] == "test"
        assert d["weight"] == 1.0
        assert d["body"] == []
        assert "confidence" not in d  # omitted when 1.0

    def test_to_dict_with_confidence(self):
        bb = BranchBody(confidence=0.5)
        d = bb.to_dict()
        assert d["confidence"] == 0.5

    def test_to_dict_with_meta(self):
        bb = BranchBody(meta={"key": "val"})
        d = bb.to_dict()
        assert d["meta"] == {"key": "val"}

    def test_from_dict(self):
        d = {"label": "A", "weight": 0.8, "body": [{"op": "tell"}], "confidence": 0.7, "meta": {"x": 1}}
        bb = BranchBody.from_dict(d)
        assert bb.label == "A"
        assert bb.weight == 0.8
        assert bb.body == [{"op": "tell"}]
        assert bb.confidence == 0.7
        assert bb.meta == {"x": 1}

    def test_from_dict_missing_fields(self):
        bb = BranchBody.from_dict({})
        assert bb.label == ""
        assert bb.weight == 1.0
        assert bb.body == []
        assert bb.confidence == 1.0

    def test_round_trip(self):
        original = BranchBody(label="fast", weight=0.8, body=[{"op": "tell", "to": "a"}], confidence=0.9)
        recovered = round_trip(original)
        assert recovered.label == "fast"
        assert recovered.weight == 0.8
        assert recovered.body == [{"op": "tell", "to": "a"}]
        assert recovered.confidence == 0.9


# ===========================================================================
# BranchPrimitive Tests
# ===========================================================================

class TestBranchPrimitive:
    def test_default_values(self):
        b = BranchPrimitive()
        assert b.strategy == "parallel"
        assert b.branches == []
        assert b.confidence == 1.0
        assert b.id != ""
        assert b.merge_strategy == "weighted_confidence"
        assert b.merge_timeout_ms == 30000
        assert b.merge_fallback == "first_complete"

    def test_confidence_clamped_high(self):
        b = BranchPrimitive(confidence=1.5)
        assert b.confidence == 1.0

    def test_confidence_clamped_low(self):
        b = BranchPrimitive(confidence=-0.5)
        assert b.confidence == 0.0

    def test_json_round_trip(self):
        original = BranchPrimitive(
            strategy="competitive",
            branches=[
                BranchBody(label="fast", weight=0.8, body=[{"op": "tell", "to": "a"}]),
                BranchBody(label="slow", body=[{"op": "ask", "from": "b"}]),
            ],
            merge_strategy="vote",
        )
        recovered = round_trip(original)
        assert recovered.strategy == "competitive"
        assert len(recovered.branches) == 2
        assert recovered.branches[0].label == "fast"
        assert recovered.branches[0].weight == 0.8
        assert recovered.branches[0].body == [{"op": "tell", "to": "a"}]
        assert recovered.branches[1].body == [{"op": "ask", "from": "b"}]
        assert recovered.merge_strategy == "vote"

    def test_schema_version_in_dict(self):
        b = BranchPrimitive()
        d = b.to_dict()
        assert d["$schema"] == "flux.a2a.branch/v1"
        assert d["op"] == "branch"

    def test_dict_construction(self):
        data = {
            "op": "branch",
            "strategy": "parallel",
            "branches": [{"label": "A", "body": []}],
            "merge": {"strategy": "consensus", "timeout_ms": 10000},
        }
        b = BranchPrimitive.from_dict(data)
        assert b.strategy == "parallel"
        assert len(b.branches) == 1
        assert b.merge_strategy == "consensus"
        assert b.merge_timeout_ms == 10000

    def test_branch_body_weight_clamped(self):
        bb = BranchBody(weight=2.0)
        assert bb.weight == 1.0

    def test_branch_with_meta(self):
        b = BranchPrimitive(meta={"author": "test"})
        d = b.to_dict()
        assert d["meta"] == {"author": "test"}

    def test_branch_with_low_confidence(self):
        b = BranchPrimitive(confidence=0.3)
        d = b.to_dict()
        assert d["confidence"] == 0.3

    def test_branch_with_all_strategies(self):
        for strategy in ["parallel", "sequential", "competitive"]:
            b = BranchPrimitive(strategy=strategy)
            assert b.strategy == strategy

    def test_merge_config_in_dict(self):
        b = BranchPrimitive(merge_timeout_ms=5000, merge_fallback="consensus")
        d = b.to_dict()
        assert d["merge"]["timeout_ms"] == 5000
        assert d["merge"]["fallback"] == "consensus"


# ===========================================================================
# ForkMutation and ForkInherit Tests
# ===========================================================================

class TestForkMutation:
    def test_default_values(self):
        m = ForkMutation()
        assert m.type == "strategy"
        assert m.changes == {}

    def test_to_dict(self):
        m = ForkMutation(type="context", changes={"risk": "high"})
        d = m.to_dict()
        assert d["type"] == "context"
        assert d["changes"] == {"risk": "high"}

    def test_from_dict(self):
        m = ForkMutation.from_dict({"type": "prompt", "changes": {"temp": 0.8}})
        assert m.type == "prompt"
        assert m.changes == {"temp": 0.8}

    def test_round_trip(self):
        original = ForkMutation(type="capability", changes={"new_skill": "dreamer"})
        recovered = round_trip(original)
        assert recovered.type == "capability"
        assert recovered.changes == {"new_skill": "dreamer"}


class TestForkInherit:
    def test_default_values(self):
        i = ForkInherit()
        assert i.state == []
        assert i.context is True
        assert i.trust_graph is False
        assert i.message_history is False

    def test_to_dict(self):
        i = ForkInherit(state=["x", "y"], trust_graph=True)
        d = i.to_dict()
        assert d["state"] == ["x", "y"]
        assert d["context"] is True
        assert d["trust_graph"] is True
        assert d["message_history"] is False

    def test_from_dict(self):
        i = ForkInherit.from_dict({"state": ["a"], "context": False})
        assert i.state == ["a"]
        assert i.context is False

    def test_round_trip(self):
        original = ForkInherit(state=["x", "y"], trust_graph=True, message_history=True)
        recovered = round_trip(original)
        assert recovered.state == ["x", "y"]
        assert recovered.trust_graph is True
        assert recovered.message_history is True


# ===========================================================================
# ForkPrimitive Tests
# ===========================================================================

class TestForkPrimitive:
    def test_default_values(self):
        f = ForkPrimitive()
        assert f.on_complete == "merge"
        assert f.conflict_mode == "negotiate"
        assert f.mutations == []
        assert f.inherit.context is True
        assert f.inherit.trust_graph is False

    def test_json_round_trip(self):
        original = ForkPrimitive(
            inherit=ForkInherit(state=["x", "y"], trust_graph=True),
            mutations=[ForkMutation(type="strategy", changes={"risk": "high"})],
            on_complete="collect",
        )
        recovered = round_trip(original)
        assert recovered.inherit.state == ["x", "y"]
        assert recovered.inherit.trust_graph is True
        assert len(recovered.mutations) == 1
        assert recovered.mutations[0].type == "strategy"
        assert recovered.mutations[0].changes == {"risk": "high"}
        assert recovered.on_complete == "collect"

    def test_schema_version(self):
        f = ForkPrimitive()
        assert f.to_dict()["$schema"] == "flux.a2a.fork/v1"

    def test_fork_with_dict_inherit(self):
        f = ForkPrimitive(inherit={"state": ["a"]})
        assert f.inherit.state == ["a"]

    def test_fork_with_dict_mutations(self):
        f = ForkPrimitive(mutations=[{"type": "prompt", "changes": {}}])
        assert f.mutations[0].type == "prompt"

    def test_fork_confidence_clamped(self):
        f = ForkPrimitive(confidence=2.0)
        assert f.confidence == 1.0

    def test_fork_all_on_complete_values(self):
        for val in ["collect", "discard", "signal", "merge"]:
            f = ForkPrimitive(on_complete=val)
            assert f.on_complete == val

    def test_fork_all_conflict_modes(self):
        for val in ["parent_wins", "child_wins", "negotiate"]:
            f = ForkPrimitive(conflict_mode=val)
            assert f.conflict_mode == val


# ===========================================================================
# CoIteratePrimitive Tests
# ===========================================================================

class TestCoIteratePrimitive:
    def test_default_values(self):
        c = CoIteratePrimitive()
        assert c.agents == []
        assert c.shared_state_mode == "merge"
        assert c.convergence_threshold == 0.95

    def test_json_round_trip(self):
        original = CoIteratePrimitive(
            agents=["oracle1", "superz"],
            shared_state_mode="partitioned",
            convergence_metric="confidence_delta",
            convergence_threshold=0.99,
        )
        recovered = round_trip(original)
        assert recovered.agents == ["oracle1", "superz"]
        assert recovered.shared_state_mode == "partitioned"
        assert recovered.convergence_metric == "confidence_delta"
        assert recovered.convergence_threshold == 0.99

    def test_schema_version(self):
        c = CoIteratePrimitive()
        assert c.to_dict()["$schema"] == "flux.a2a.co_iterate/v1"

    def test_convergence_threshold_clamped(self):
        c = CoIteratePrimitive(convergence_threshold=1.5)
        assert c.convergence_threshold == 1.0

    def test_shared_state_modes(self):
        for val in ["conflict", "merge", "partitioned", "isolated"]:
            c = CoIteratePrimitive(shared_state_mode=val)
            assert c.shared_state_mode == val

    def test_co_iterate_with_meta(self):
        c = CoIteratePrimitive(meta={"notes": "test"})
        d = c.to_dict()
        assert d["meta"] == {"notes": "test"}

    def test_from_dict_with_convergence(self):
        data = {"convergence": {"metric": "agreement", "threshold": 0.5}}
        c = CoIteratePrimitive.from_dict(data)
        assert c.convergence_metric == "agreement"
        assert c.convergence_threshold == 0.5


# ===========================================================================
# Participant Tests
# ===========================================================================

class TestParticipant:
    def test_default_values(self):
        p = Participant()
        assert p.agent == ""
        assert p.stance == "neutral"
        assert p.role == ""

    def test_to_dict_minimal(self):
        p = Participant(agent="oracle1")
        d = p.to_dict()
        assert d["agent"] == "oracle1"
        assert d["stance"] == "neutral"
        assert "role" not in d

    def test_to_dict_with_role(self):
        p = Participant(agent="oracle1", role="moderator")
        d = p.to_dict()
        assert d["role"] == "moderator"

    def test_from_dict(self):
        p = Participant.from_dict({"agent": "test", "stance": "con", "role": "devil"})
        assert p.agent == "test"
        assert p.stance == "con"
        assert p.role == "devil"

    def test_round_trip(self):
        original = Participant(agent="superz", stance="pro", role="moderator")
        recovered = round_trip(original)
        assert recovered.agent == "superz"
        assert recovered.stance == "pro"
        assert recovered.role == "moderator"

    def test_all_stances(self):
        for stance in ["pro", "con", "neutral", "devil's_advocate", "moderator"]:
            p = Participant(stance=stance)
            assert p.stance == stance


# ===========================================================================
# DiscussPrimitive Tests
# ===========================================================================

class TestDiscussPrimitive:
    def test_default_values(self):
        d = DiscussPrimitive()
        assert d.format == "peer_review"
        assert d.turn_order == "round_robin"
        assert d.until_condition == "consensus"
        assert d.max_rounds == 5

    def test_json_round_trip(self):
        original = DiscussPrimitive(
            format="debate",
            topic="Binary vs JSON messages",
            participants=[
                Participant(agent="oracle1", stance="pro", role="moderator"),
                Participant(agent="superz", stance="neutral"),
            ],
            max_rounds=10,
        )
        recovered = round_trip(original)
        assert recovered.format == "debate"
        assert recovered.topic == "Binary vs JSON messages"
        assert len(recovered.participants) == 2
        assert recovered.participants[0].agent == "oracle1"
        assert recovered.participants[0].stance == "pro"
        assert recovered.participants[0].role == "moderator"
        assert recovered.participants[1].stance == "neutral"
        assert recovered.max_rounds == 10

    def test_schema_version(self):
        d = DiscussPrimitive(topic="test")
        assert d.to_dict()["$schema"] == "flux.a2a.discuss/v1"

    def test_from_dict_with_until(self):
        data = {"until": {"condition": "timeout", "max_rounds": 20}}
        d = DiscussPrimitive.from_dict(data)
        assert d.until_condition == "timeout"
        assert d.max_rounds == 20

    def test_discuss_with_dict_participants(self):
        d = DiscussPrimitive(participants=[{"agent": "a", "stance": "pro"}])
        assert d.participants[0].agent == "a"

    def test_discuss_all_formats(self):
        for fmt in ["debate", "brainstorm", "review", "negotiate", "peer_review"]:
            d = DiscussPrimitive(format=fmt)
            assert d.format == fmt


# ===========================================================================
# SynthesisSource Tests
# ===========================================================================

class TestSynthesisSource:
    def test_default_values(self):
        s = SynthesisSource()
        assert s.type == "variable"
        assert s.ref == ""

    def test_to_dict(self):
        s = SynthesisSource(type="branch_result", ref="exploration")
        d = s.to_dict()
        assert d["type"] == "branch_result"
        assert d["ref"] == "exploration"

    def test_from_dict(self):
        s = SynthesisSource.from_dict({"type": "external", "ref": "human"})
        assert s.type == "external"
        assert s.ref == "human"

    def test_round_trip(self):
        original = SynthesisSource(type="fork_result", ref="child_1")
        recovered = round_trip(original)
        assert recovered.type == "fork_result"
        assert recovered.ref == "child_1"


# ===========================================================================
# SynthesizePrimitive Tests
# ===========================================================================

class TestSynthesizePrimitive:
    def test_default_values(self):
        s = SynthesizePrimitive()
        assert s.method == "map_reduce"
        assert s.output_type == "decision"
        assert s.confidence_mode == "propagate"

    def test_json_round_trip(self):
        original = SynthesizePrimitive(
            method="ensemble",
            sources=[
                SynthesisSource(type="branch_result", ref="exploration"),
                SynthesisSource(type="external", ref="human_feedback"),
            ],
            output_type="summary",
        )
        recovered = round_trip(original)
        assert recovered.method == "ensemble"
        assert len(recovered.sources) == 2
        assert recovered.sources[0].type == "branch_result"
        assert recovered.sources[1].ref == "human_feedback"
        assert recovered.output_type == "summary"

    def test_schema_version(self):
        s = SynthesizePrimitive()
        assert s.to_dict()["$schema"] == "flux.a2a.synthesize/v1"

    def test_with_dict_sources(self):
        s = SynthesizePrimitive(sources=[{"type": "variable", "ref": "x"}])
        assert s.sources[0].type == "variable"

    def test_all_methods(self):
        for method in ["map_reduce", "ensemble", "chain", "vote", "weighted_merge", "best_effort"]:
            s = SynthesizePrimitive(method=method)
            assert s.method == method

    def test_all_output_types(self):
        for ot in ["code", "spec", "question", "decision", "summary", "value"]:
            s = SynthesizePrimitive(output_type=ot)
            assert s.output_type == ot

    def test_all_confidence_modes(self):
        for cm in ["propagate", "min", "max", "average"]:
            s = SynthesizePrimitive(confidence_mode=cm)
            assert s.confidence_mode == cm


# ===========================================================================
# ReflectPrimitive Tests
# ===========================================================================

class TestReflectPrimitive:
    def test_default_values(self):
        r = ReflectPrimitive()
        assert r.target == "strategy"
        assert r.method == "introspection"
        assert r.output == "adjustment"

    def test_json_round_trip(self):
        original = ReflectPrimitive(
            target="progress",
            method="benchmark",
            output="question",
            confidence=0.7,
        )
        recovered = round_trip(original)
        assert recovered.target == "progress"
        assert recovered.method == "benchmark"
        assert recovered.output == "question"
        assert recovered.confidence == 0.7

    def test_schema_version(self):
        r = ReflectPrimitive()
        assert r.to_dict()["$schema"] == "flux.a2a.reflect/v1"

    def test_all_targets(self):
        for target in ["strategy", "progress", "uncertainty", "confidence", "all"]:
            r = ReflectPrimitive(target=target)
            assert r.target == target

    def test_all_methods(self):
        for method in ["introspection", "benchmark", "comparison", "statistical"]:
            r = ReflectPrimitive(method=method)
            assert r.method == method

    def test_all_outputs(self):
        for output in ["adjustment", "question", "branch", "log", "signal"]:
            r = ReflectPrimitive(output=output)
            assert r.output == output

    def test_reflect_with_meta(self):
        r = ReflectPrimitive(meta={"iteration": 5})
        d = r.to_dict()
        assert d["meta"] == {"iteration": 5}


# ===========================================================================
# Registry Tests
# ===========================================================================

class TestRegistry:
    def test_all_primitives_registered(self):
        assert len(PRIMITIVE_REGISTRY) == 6
        assert "branch" in PRIMITIVE_REGISTRY
        assert "fork" in PRIMITIVE_REGISTRY
        assert "co_iterate" in PRIMITIVE_REGISTRY
        assert "discuss" in PRIMITIVE_REGISTRY
        assert "synthesize" in PRIMITIVE_REGISTRY
        assert "reflect" in PRIMITIVE_REGISTRY

    def test_parse_branch(self):
        result = parse_primitive({"op": "branch", "branches": []})
        assert isinstance(result, BranchPrimitive)

    def test_parse_fork(self):
        result = parse_primitive({"op": "fork"})
        assert isinstance(result, ForkPrimitive)

    def test_parse_co_iterate(self):
        result = parse_primitive({"op": "co_iterate", "agents": ["a"]})
        assert isinstance(result, CoIteratePrimitive)

    def test_parse_discuss(self):
        result = parse_primitive({"op": "discuss", "topic": "test"})
        assert isinstance(result, DiscussPrimitive)

    def test_parse_synthesize(self):
        result = parse_primitive({"op": "synthesize", "sources": []})
        assert isinstance(result, SynthesizePrimitive)

    def test_parse_reflect(self):
        result = parse_primitive({"op": "reflect"})
        assert isinstance(result, ReflectPrimitive)

    def test_parse_unknown_returns_none(self):
        result = parse_primitive({"op": "let", "name": "x", "value": 1})
        assert result is None

    def test_parse_no_op_returns_none(self):
        result = parse_primitive({"not": "an op"})
        assert result is None

    def test_parse_empty_dict_returns_none(self):
        result = parse_primitive({})
        assert result is None

    def test_parse_none_raises_error(self):
        """parse_primitive(None) should raise AttributeError since None has no .get()."""
        with pytest.raises(AttributeError):
            parse_primitive(None)


# ===========================================================================
# Integration / Cross-Primitive Tests
# ===========================================================================

class TestIntegration:
    """Test primitives working together and edge cases."""

    def test_branch_with_explicit_id(self):
        b = BranchPrimitive(id="my-branch-id")
        assert b.id == "my-branch-id"

    def test_fork_with_explicit_id(self):
        f = ForkPrimitive(id="my-fork-id")
        assert f.id == "my-fork-id"

    def test_all_primitives_have_ids(self):
        """All primitives should auto-generate UUIDs."""
        for cls in [BranchPrimitive, ForkPrimitive, CoIteratePrimitive,
                    DiscussPrimitive, SynthesizePrimitive, ReflectPrimitive]:
            p = cls()
            assert len(p.id) > 0

    def test_all_primitives_have_schema_field(self):
        """All primitives should have a $schema field in to_dict output."""
        schemas = {
            "branch": "flux.a2a.branch/v1",
            "fork": "flux.a2a.fork/v1",
            "co_iterate": "flux.a2a.co_iterate/v1",
            "discuss": "flux.a2a.discuss/v1",
            "synthesize": "flux.a2a.synthesize/v1",
            "reflect": "flux.a2a.reflect/v1",
        }
        for op, expected_schema in schemas.items():
            p = parse_primitive({"op": op})
            d = p.to_dict()
            assert d["$schema"] == expected_schema

    def test_all_primitives_have_op_field(self):
        """All primitives should have 'op' field in to_dict output."""
        for op_name in PRIMITIVE_REGISTRY:
            p = parse_primitive({"op": op_name})
            d = p.to_dict()
            assert d["op"] == op_name

    def test_all_primitives_json_serializable(self):
        """All primitives should produce JSON-serializable output."""
        for op_name in PRIMITIVE_REGISTRY:
            p = parse_primitive({"op": op_name})
            d = p.to_dict()
            json_str = json.dumps(d)
            assert isinstance(json_str, str)
            assert len(json_str) > 0

    def test_complex_branch(self):
        """Test a branch with 5 sub-paths and all merge options."""
        b = BranchPrimitive(
            strategy="parallel",
            branches=[
                BranchBody(label=f"path_{i}", weight=i * 0.2, body=[{"step": i}])
                for i in range(5)
            ],
            merge_strategy="weighted_confidence",
            merge_timeout_ms=60000,
            merge_fallback="all",
            confidence=0.8,
            meta={"source": "integration_test"},
        )
        d = b.to_dict()
        assert len(d["branches"]) == 5
        recovered = round_trip(b)
        assert len(recovered.branches) == 5
        assert recovered.branches[4].label == "path_4"

    def test_complex_fork(self):
        """Test a fork with multiple mutations and full inheritance."""
        f = ForkPrimitive(
            inherit=ForkInherit(
                state=["knowledge", "context", "trust_level"],
                context=True,
                trust_graph=True,
                message_history=True,
            ),
            mutations=[
                ForkMutation(type="strategy", changes={"risk_tolerance": "high"}),
                ForkMutation(type="prompt", changes={"system_msg": "new"}),
                ForkMutation(type="capability", changes={"new_skills": ["dreamer"]}),
            ],
            on_complete="merge",
            conflict_mode="negotiate",
            confidence=0.9,
        )
        d = f.to_dict()
        assert len(d["mutations"]) == 3
        assert d["inherit"]["state"] == ["knowledge", "context", "trust_level"]
        recovered = round_trip(f)
        assert len(recovered.mutations) == 3

    def test_complex_discuss(self):
        """Test a discussion with 6 participants."""
        roles = ["architect", "critic", "pragmatist", "visionary", "historian", "contrarian"]
        d = DiscussPrimitive(
            format="debate",
            topic="Should FLUX use variable-width opcodes?",
            participants=[Participant(agent=r, stance="neutral") for r in roles],
            max_rounds=10,
            confidence=0.85,
        )
        recovered = round_trip(d)
        assert len(recovered.participants) == 6

    def test_primitives_to_dict_no_crash_on_empty(self):
        """All primitives should serialize cleanly with defaults."""
        for cls in [BranchPrimitive, ForkPrimitive, CoIteratePrimitive,
                    DiscussPrimitive, SynthesizePrimitive, ReflectPrimitive]:
            p = cls()
            d = p.to_dict()
            assert isinstance(d, dict)
            assert len(d) > 0
