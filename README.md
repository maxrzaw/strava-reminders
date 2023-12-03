# Strava Reminders

This app aims to support sending push notifications based on two scenarios:

1. An activity is uploaded and the gear is not set after a time limit.
2. An activity is uploaded and the gear is still primary after a time limit.

Can specify per activity configuration for time limit and notification scenarios.

## TODO

- Subscribe to webhooks from the strava api
- I'll probably need to do something with the refresh token
- Business logic in the server
- Authorization for content specific to the logged in user
  - Hopefully I can do this with middleware
- Lightweight notification solution
  - Either printing to console or I could use [ntfy](https://ntfy.sh)
- UI for changing settings
- Make UI Mobile friendly and make it a PWA
- Add real push notifications to the UI

## DONE

- Linking a user to their Strava account and saving the tokens to the database
- Authentication and Authorization of some sort. A user can only look at their stuff, not others
  - Can I piggyback off of Strava?
    - Yes!
