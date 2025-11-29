# Development Guide

The MySQL image mainly serves to **initialize database data**, with the core being the `gvb.sql` file.

This SQL file runs automatically when the MySQL container starts.

To change initial data, edit `gvb.sql`.

> If the project has already run once, delete existing data under `start/gvb` first (remember to back up).

Then rerun the one-click script: `./bootstrap.sh`