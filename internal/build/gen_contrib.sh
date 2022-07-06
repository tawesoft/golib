git log --date=format:"%Y" --pretty=format:"%aN <%aE> (%ad)" \
    | sort -u > CONTRIBUTORS.txt
