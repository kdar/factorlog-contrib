gelf
====

This formatter outputs [GELF v1.1](http://graylog2.org/gelf#specs)

This is used in our production system to output GELF formatted logs. The features built into it (such as logging http.Request), was because that's what we needed. If you have different requirements, feel free to copy the code into your repository and modify it.

I may make the gelf package more generic so others can do what they want. If you have an idea of how to do this with a good API, let me know.
