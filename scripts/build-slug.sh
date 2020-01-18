PROJ_DIR=`pwd`
pushd $TEMP
cp -r "${PROJ_DIR}/bin" ./app
cp "${PROJ_DIR}/scripts/start.sh" ./app
tar czfv slug.tgz ./app
rm -r app
popd
mv "$TEMP/slug.tgz" ./