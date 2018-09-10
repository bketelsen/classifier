## Classifier

## Assumptions

- You have scanned statements in a directory called "Unfiled".  The scans have already had OCR processed and the OCR text has been added to the PDF. 


### Getting Started
Place a file called `classifier.toml` in your home directory.  Use the included example as a starting point.
Modify this file to have unique keywords for each "vendor".  Account numbers, payment addresses, vendor name are all good choices for unique keywords.

Install pdf2txt from `https://github.com/euske/pdfminer/blob/master/tools/pdf2txt.py`

`go install`

`classifier -d /Users/you/Documents`

`classifier` will read from `/Users/you/Documents/Unfiled` and file into
`/Users/you/Documents/Filed/{vendor}/year/month`

### This WILL EAT your files.  
Make backups.  No guarantees, you've been warned.