FROM alpine

COPY gke-preemptible-notifier /usr/bin/

ENTRYPOINT ["gke-preemptible-notifier", "watch"]
