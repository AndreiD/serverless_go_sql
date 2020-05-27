gcloud functions deploy serverlesskata --env-vars-file .env.yml  \
--runtime go113 --trigger-http --allow-unauthenticated --entry-point Register \
--memory 256MB --region europe-west3 --max-instances 1 \
--source=.