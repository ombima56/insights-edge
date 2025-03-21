// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Marketplace {
    struct User {
        string username;
        string userType;
        uint256 balance;
    }

    struct Insight {
        string industry;
        string insightType;
        string content;
        address author;
        uint256 price;
    }

    mapping(address => User) public users;
    Insight[] public insights;
    mapping(address => bool) public subscribers;
    address public owner;

    event InsightCreated(address indexed author, string industry, string insightType, uint256 price);
    event InsightPurchased(address indexed buyer, uint256 insightId);
    event Subscribed(address indexed user, string plan);

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    function registerUser(string memory _username, string memory _userType) external {
        users[msg.sender] = User(_username, _userType, 100); // Default balance 100 BPT
    }

    function createInsight(string memory _industry, string memory _type, string memory _content, uint256 _price) external {
        insights.push(Insight(_industry, _type, _content, msg.sender, _price));
        emit InsightCreated(msg.sender, _industry, _type, _price);
    }

    function purchaseInsight(uint256 _insightId) external {
        require(_insightId < insights.length, "Invalid insight ID");
        Insight storage insight = insights[_insightId];
        require(users[msg.sender].balance >= insight.price, "Insufficient balance");
        
        users[msg.sender].balance -= insight.price;
        users[insight.author].balance += insight.price;
        
        emit InsightPurchased(msg.sender, _insightId);
    }

    function subscribe(string memory _plan) external {
        require(!subscribers[msg.sender], "Already subscribed");
        subscribers[msg.sender] = true;
        emit Subscribed(msg.sender, _plan);
    }

    function getUser(address _user) external view returns (string memory, string memory, uint256) {
        User memory user = users[_user];
        return (user.username, user.userType, user.balance);
    }

    function getInsight(uint256 _insightId) external view returns (string memory, string memory, string memory, address, uint256) {
        require(_insightId < insights.length, "Invalid insight ID");
        Insight memory insight = insights[_insightId];
        return (insight.industry, insight.insightType, insight.content, insight.author, insight.price);
    }
}
