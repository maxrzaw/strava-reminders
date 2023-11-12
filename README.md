# Strava Reminders

This app aims to support sending push notifications based on two scenarios:
1. An activity is uploaded and the gear is not set after a time limit.
2. An activity is uploaded and the gear is still primary after a time limit.

Can specify per activity configuration for time limit and notification scenarios.


## TODO

- Business logic in the server
- Authentication and Authorization of some sort
  - A user can only look at their stuff, not others
  - Can I piggyback off of strava?
    - I don't think so :(
    - Can probably use a JWT token in a cookie and write and validate it using go middleware
    - This [Reddit post](https://www.reddit.com/r/htmx/comments/11jwxeq/authentication_for_htmx_app/) has some info
- Lightweight notification solution
  - Either printing to console or I could use [ntfy](https://ntfy.sh)
- UI for changing settings
- Make UI Mobile friendly and make it a PWA
- Add real push notifications to the UI
