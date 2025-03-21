// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Dashboard {
    struct User {
        string username;
        string userType;
        address wallet;
        uint256 balance;
    }
    
    struct Insight {
        string industry;
        string insightType;
        string content;
        address author;
    }
    
    struct Subscription {
        string plan;
        uint256 expiresAt;
    }
    
    mapping(address => User) public users;
    mapping(address => Insight[]) public insights;
    mapping(address => Subscription) public subscriptions;
    
    event UserRegistered(address user, string username, string userType);
    event InsightCreated(address indexed user, string industry, string insightType);
    event Subscribed(address indexed user, string plan, uint256 expiresAt);
    
    function registerUser(string memory _username, string memory _userType) public {
        users[msg.sender] = User(_username, _userType, msg.sender, 100);
        emit UserRegistered(msg.sender, _username, _userType);
    }
    
    function createInsight(string memory _industry, string memory _type, string memory _content) public {
        insights[msg.sender].push(Insight(_industry, _type, _content, msg.sender));
        emit InsightCreated(msg.sender, _industry, _type);
    }
    
    function subscribe(string memory _plan) public {
        uint256 duration = 30 days;
        subscriptions[msg.sender] = Subscription(_plan, block.timestamp + duration);
        emit Subscribed(msg.sender, _plan, block.timestamp + duration);
    }
    
    function getUser(address _user) public view returns (User memory) {
        return users[_user];
    }
    
    function getInsights(address _user) public view returns (Insight[] memory) {
        return insights[_user];
    }
    
    function getSubscription(address _user) public view returns (Subscription memory) {
        return subscriptions[_user];
    }
}
