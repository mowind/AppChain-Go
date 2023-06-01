// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.12;

interface ICheckpointSigAggregator {
    event CheckpointSigAggregated(
        address indexed proposer,
        uint256 start,
        uint256 end,
        bytes32 root,
        uint32[] signedValidators,
        bytes signature
    );

    struct Checkpoint {
        address proposer;
        uint256 start;
        uint256 end;
        bytes32 rootHash;
        bytes32 accountHash;
        uint256 chainId;
        uint32[] current;
        uint32[] rewards;
        uint32[] slashing;
    }

    function propose(
        Checkpoint calldata cp,
        uint32 validatorId,
        bytes calldata signature
    ) external;

    function confirm(address proposer, bytes32 root) external;

    function latestCheckpoint() external view returns (Checkpoint memory cp);
}
