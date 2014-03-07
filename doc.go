/*
This package consists of necessary parts of building and querying an activity net.

An activity net is made of extracted activities, POIs, original messages and users.
It can be used to model/predict user behavior, discover new POIs, or help NLP in turn.

The workflow consists of the following parts.

  1. Process the raw data into the form required by the importer. There can be multiple importers.
  2. Use the importer to generate the universal intermediate representations of source data.
  3. Build the activity net from the source data. Revise this model using human labeled data.
  4. Query the activity net.

With this workflow, we can integrate different data sources into one.


This package can answer the following questions:

  1. Are my extracted activities real human activities? (Classification)
  2. How does my extracted activities look like? (Visualization)
  2. Is there any new activities recently? (Trend)
  3. What's the relations between activities? (Relation)

Note that this package do not focus on the recall/precision of the activity extraction.
It's rather a tool for ultilizing the output of such process.
*/
package actnet
