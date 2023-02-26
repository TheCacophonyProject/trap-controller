#!/bin/bash
systemctl daemon-reload
systemctl enable trap-controller.service
systemctl restart trap-controller.service