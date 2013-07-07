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

	$> mongo power_info
	>  db.info.ensureIndex( { "Time": 1, "Key": 1 } )
	$> mongoimport --db power_info -c info --upsertFields Time,Key --file ./tmp/power_info.json.log

I guess you could just pipe the info straight to mongo as well. Just giving ideas.

This was a hasty hack to collect stats from a new battery I bought for an old 
laptop. I wanted to be collecting stats from the time I put the battery on,
so I might see when it starts to go south. :-)

Compile/Install
---------------

Have Google Go lang installed, and run:

	go get github.com/vbatts/power_info


Bugs / Ideas
------------

Feel free to open an issue, or submit a pull request.

