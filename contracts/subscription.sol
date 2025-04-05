// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SubscriptionMarket {
    enum Plan { None, Monthly, Quarterly, Yearly }

    struct Subscription {
        Plan plan;
        uint256 expiry;
    }

    mapping(address => Subscription) public subscriptions;

    event Subscribed(address indexed user, Plan plan, uint256 expiry);

    uint256 public constant MONTHLY_PRICE = 0.01 ether;
    uint256 public constant QUARTERLY_PRICE = 0.025 ether;
    uint256 public constant YEARLY_PRICE = 0.08 ether;

    uint256 public constant MONTH = 30 days;
    uint256 public constant QUARTER = 90 days;
    uint256 public constant YEAR = 365 days;

    function subscribe(Plan _plan) external payable {
        require(_plan != Plan.None, "Invalid plan selected");

        uint256 price = getPrice(_plan);
        uint256 duration = getDuration(_plan);
        require(msg.value == price, "Incorrect ETH sent for selected plan");

        uint256 currentExpiry = subscriptions[msg.sender].expiry;
        uint256 newExpiry = block.timestamp > currentExpiry
            ? block.timestamp + duration
            : currentExpiry + duration;

        subscriptions[msg.sender] = Subscription({
            plan: _plan,
            expiry: newExpiry
        });

        emit Subscribed(msg.sender, _plan, newExpiry);
    }

    function getPrice(Plan _plan) public pure returns (uint256) {
        if (_plan == Plan.Monthly) return MONTHLY_PRICE;
        if (_plan == Plan.Quarterly) return QUARTERLY_PRICE;
        if (_plan == Plan.Yearly) return YEARLY_PRICE;
        return 0;
    }

    function getDuration(Plan _plan) public pure returns (uint256) {
        if (_plan == Plan.Monthly) return MONTH;
        if (_plan == Plan.Quarterly) return QUARTER;
        if (_plan == Plan.Yearly) return YEAR;
        return 0;
    }

    function isActive(address user) public view returns (bool) {
        return subscriptions[user].expiry > block.timestamp;
    }

    function getSubscription(address user)
        external
        view
        returns (Plan plan, uint256 expiry, bool active)
    {
        Subscription memory sub = subscriptions[user];
        return (sub.plan, sub.expiry, isActive(user));
    }
}
