## How to set up Google OAuth

In project i used `Google OAuth` to connect with Google Calendar service.

It is a several steps to set up the project:

1. You should follow this [guide (Go quickstart | Google Calendar)](https://developers.google.com/calendar/api/quickstart/go) until you get the file `credentials.json`. After that, you need to paste its contents into `credentials.json` file is in the main directory of the project.

2. After that you should run test file, to set up your `token.json` file.
    Your should do these command in your terminal:
    ```shell
    cd cmd
    go run main.go -test
    ```

    You will see a link in terminal, which should follow untill you can not be redirect to `localhost`.

    Terminal output (with link):
    ```
    INFO: 2024/05/27 17:42:54 /Users/dmitrykuzmin/Documents/Programming/projects/private-planpilot/cmd <nil>
    Go to the following link in your browser then type the authorization code:
    https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=197118240391-6tbnbiclm1aqnpjs84dg2ggflubjfift.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcalendar&state=state-token
    ```

    When you will direct to localhost, you  need to extract the authorization code from link:
    `http://localhost/?state=state-token&code=4/0AdLIrYfRbFs_8r4lyLyfgoeb1fUsCysCqyMqCwSqY7rACq1971KLwZna61jZZqxNwprkDw&scope=https://www.googleapis.com/auth/calendar`
    
    **Result code**: `40AdLIrYfRbFs_8r4lyLyfgoeb1fUsCysCqyMqCwSqY7rACq1971KLwZna61jZZqxNwprkD`
    Paste with code into terminal and you will see success message:
    ```
    oauth token:  4/0AdLIrYczCOu8BXyK5fXzWfThS1V21WCp4JqOLN-vPGqQ7PVt2TfHtgigGHR50sm0jvW_1w
    Saving credential file to: ../token.json
    ```