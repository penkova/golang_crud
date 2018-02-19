FROM bitnami/minideb
ADD main /
EXPOSE 80
CMD ["/main"]