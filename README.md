# Routini
Intelligent scheduler for any task

---

## API End points
### User 
- Signup
- Login
### Calender Entry
- Create
- Read/Get Entry
- Edit 
- Delete
### GUI 
- Get Month
- Get Week
- Get Day
- Read/Get Entry
---
## DB 
### User
- ID
- Email
- Name
- PasswordHash
- MonthIDs [ ] ?

### Month
- ID
- UserID
- Year
- MonthNumber
- Weeks [  ]
  - WeekNumber
  - Days [  ]
    - DayNumber
    - Date
    - Entry [ ]
      - Name
      - Tags [ ]
      - StartTime
      - EndTime
      - DueDate?
      - Repeat 
        - Repeat Every N Days
        - On [X] Days
        - Ends On Z Date || After Z Occurances
      - Location {Adress}
      - MeetingLink
      - Reminder [ Times ]
      - Guests [ UserIDs/emails ]
      - Notes
      - Prority
      - Permanent/Movable
---
### TODO
TODO TODO TODOTODOTODOTODOTODOOOOO TODOTODODODO