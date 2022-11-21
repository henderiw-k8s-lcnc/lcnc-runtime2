# lcnc-runtime2
 
block types:
  - for
  - own
  - watch
  - vars
  - resources
  - services

## options for lcnc syntax

option1:

provide all the data to all functions
- pro: 
  - runtime simple
- con: 
  - heap could become to big

option2:

provide query interface from fn to runtime
- pro: 
  - heap usage is more efficient
- con:
  - runtime complexity

option3:

provide a query during runtime
- pro:
  - heap
- con:
  - runtime functionality: dag, loop