# Requirements

### Where used functionality.
- Return all rules that include a specified host/network/range/port/group. 
- Specific match - Only returns for rules that match the search exactly.
- Generic match - Return all rules that contain the search i.e. searching for 192.168.1.1 will also return rules that have 192.168.1.0/24.
- Return all the rules that include an object that falls within the search range/network.

### Access checker
- User can check if a connection is allowed. Works out the route it would take and tests against each firewall in the path.

### Group Queries
- Allow user to query group info.
- Return the members of the search group.
- Allow the user to find every group a particular object is a member.
