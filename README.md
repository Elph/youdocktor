# youdocktor

this tool is first and foremost a way to play with golang with somthing more or less real. The tool will syncronize the time doctor spent time 
in tasks with the youtrack associated tasks.. 
It relies on some convention over the tasks names in TimeDoctor

Task name should begin with the format [YYY-XXX] being the youtrack card ID [PROJECT-NUMBER]
if found a time for a given Youtrack card it will update/create a workitem entry that starts with a specific text

```
(youdocktor:<timedoctor taskid>:date) <text of the task in the timedoctor task>
```

So multiple entries will be entered if the task has been worked on in different days. In a specific day it will sum all the spent time.


## How to make it work

First of all you need to create a copy of the `config.example.json` and rename it to `config.json`, and fill the contents with your own TimeDoctor / Youtrack Tokens.

For TimeDoctor, you'll need to create your own App, so you're the only one who will access your data
For Youtrack, you just need to create a token from your profile page

Once done you'll be able to call the program with some flags

* from: starting date you want to start requesting and update times in YYYY-MM-DD format
* to: ending date you want to start requesting and update times in YYYY-MM-DD format
* dry-run: Boolean flag, it will run the tool without updating anything on youtrack
* verbose: Will set the logs to `Trace` instead the default `Warn`
* summary: Will output a table with the list of work logs and associated Youtrack cards 

The summary output will be something like 

```
|    DATE    | STATUS | IDENTIFIER |              NAME              |  SPENT   |
+------------+--------+------------+--------------------------------+----------+
| 2019-03-01 | +      | XX-915     | [XX-915] task 1                | 3h27m49s |
| 2019-03-01 | ?      | n/a        | support                        | 50m36s   |
| 2019-03-01 | +      | XX-1842    | [XX-1842] support              | 39m35s   |
| 2019-03-01 | ?      | n/a        | meeting                        | 25m53s   |
| 2019-03-04 | =      | n/a        | update stuf                    | 1h49m13s |
| 2019-03-04 | +      | YY-915     | [YY-915] call                  | 1h38m12s |
| 2019-03-04 | +      | ZZ-853     | [ZZ-853] more stuff            | 1h22m55s |
```

The `Status` column defines how it will be sync to Youtrack

* `+` : a new time entry in youtrack will be created
* `>` : the entry in youtrack already exists but it will be updated
* `=` : the entry in youtrack already exists but won't be updated because the time is the same
* `?` : no match in Youtrack.

## Others

* Different tasks in TimeDoctor will create diferent time tracking entries in Youtrack even if the YouTrack card is the same
* Once added the tokens in the config file, the TimeDoctor ones should be updated once the accces token expires

