power_info
==========

Utility to collect power supply, loadavg and version on Linux OS

Overview
--------

If you just run it like:

	$> power_info


Or add it to a crontab to continue collection to a file:

	*/1 * * * * ~/bin/power_info -quiet >> ~/tmp/power_info.json.log

Read the file into mongo for playing with, or aggregating the numbers

	mongoimport --db power_info -c info --upsert --file ./tmp/power_info.json.log

I guess you could just pipe the info straight to mongo as well. Just giving ideas.

Compile/Install
---------------


