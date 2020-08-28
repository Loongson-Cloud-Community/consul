export default () => key => {
  if (key.startsWith('Status:')) {
    const [, dir] = key.split(':');
    const props = [
      'PercentageMeshChecksPassing',
      'PercentageMeshChecksWarning',
      'PercentageMeshChecksCritical',
    ];
    if (dir === 'asc') {
      props.reverse();
    }
    return function(a, b) {
      for (let i in props) {
        let prop = props[i];
        if (a[prop] === b[prop]) {
          continue;
        }
        return a[prop] > b[prop] ? -1 : 1;
      }
    };
  }
  return key;
};
