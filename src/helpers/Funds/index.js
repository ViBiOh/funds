import { buildFullTextRegex, fullTextRegexFilter } from "helpers/Search";
import { isUndefined } from "helpers/Object";

export function filterFunds(funds, filters = {}) {
  if (!Array.isArray(funds)) {
    return [];
  }

  return Object.keys(filters).reduce((previous, filter) => {
    const regex = buildFullTextRegex(String(filters[filter]));
    return previous.filter(fund => fullTextRegexFilter(fund[filter], regex));
  }, funds.slice());
}

export function orderFunds(funds, { key, descending }) {
  if (!key) {
    return funds;
  }

  const compareMultiplier = descending ? -1 : 1;

  funds.sort((o1, o2) => {
    if (isUndefined(o1, key)) {
      return -1 * compareMultiplier;
    }
    if (isUndefined(o2, key)) {
      return 1 * compareMultiplier;
    }
    if (o1[key] < o2[key]) {
      return -1 * compareMultiplier;
    }
    if (o1[key] > o2[key]) {
      return 1 * compareMultiplier;
    }
    return 0;
  });

  return funds;
}

export function aggregateFunds(funds, aggregat) {
  if (!aggregat.key) {
    return [];
  }

  const aggregate = {};
  const size = Math.min(funds.length, aggregat.size);
  for (let i = 0; i < size; i += 1) {
    if (isUndefined(aggregate, funds[i][aggregat.key])) {
      aggregate[funds[i][aggregat.key]] = 0;
    }
    aggregate[funds[i][aggregat.key]] += 1;
  }

  const aggregated = Object.keys(aggregate).map(label => ({
    label,
    count: aggregate[label]
  }));
  aggregated.sort((o1, o2) => o2.count - o1.count);

  return aggregated;
}
