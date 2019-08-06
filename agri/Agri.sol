pragma solidity >=0.4.22 <0.6.0;
contract Agri {


    bytes32[] public farmPicHash;

    function addNewPic(bytes32 picHash) public {
        farmPicHash.push(picHash);
    }

}
