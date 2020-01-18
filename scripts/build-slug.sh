PROJ_DIR=`pwd`
pushd $TEMP
cp -r "${PROJ_DIR}/bin" ./bin
cp "${PROJ_DIR}/scripts/start.sh" ./bin
tar czfv slug.tgz ./bin
rm -r bin
popd
mv "$TEMP/slug.tgz" ./