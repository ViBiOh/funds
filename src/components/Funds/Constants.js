/**
 * Columns description.
 * @type {Object}
 */
export const COLUMNS = {
  isin: {
    label: "ISIN",
    sortable: true,
    summable: false,
    filterable: true
  },
  label: {
    label: "Libellé",
    sortable: true,
    summable: false,
    filterable: true
  },
  category: {
    label: "Catégorie",
    sortable: true,
    summable: true,
    filterable: true
  },
  rating: {
    label: "Note",
    sortable: true,
    summable: true,
    filterable: true
  },
  "1m": {
    label: "1 mois",
    sortable: true,
    summable: false,
    filterable: false
  },
  "3m": {
    label: "3 mois",
    sortable: true,
    summable: false,
    filterable: false
  },
  "6m": {
    label: "6 mois",
    sortable: true,
    summable: false,
    filterable: false
  },
  "1y": {
    label: "1 an",
    sortable: true,
    summable: false,
    filterable: false
  },
  v3y: {
    label: "Volatilité",
    sortable: true,
    summable: false,
    filterable: false
  },
  score: {
    label: "Score",
    sortable: true,
    summable: false,
    filterable: false
  }
};

/**
 * Columns header label.
 * @type {Object}
 */
export const COLUMNS_HEADER = Object.entries(COLUMNS).reduce(
  (previous, [key, value]) => {
    // eslint-disable-next-line no-param-reassign
    previous[key] = value.label;
    return previous;
  },
  {}
);

/**
 * Color cycle for charts.
 * @type {Array}
 */
export const CHART_COLORS = [
  "#1f77b4",
  "#e377c2",
  "#ff7f0e",
  "#2ca02c",
  "#bcbd22",
  "#d62728",
  "#17becf",
  "#9467bd",
  "#7f7f7f",
  "#8c564b",
  "#3366cc"
];

/**
 * Default aggregate size.
 * @type {Array}
 */
export const AGGREGATE_SIZES = [25, 50, 100];

export const AGGREGAT_PARAM = "a";
export const AGGREGAT_SIZE_PARAM = "as";
export const ORDER_PARAM = "o";
export const ASCENDING_ORDER_PARAM = "ao";
export const RESERVED_PARAM = [
  AGGREGAT_PARAM,
  AGGREGAT_SIZE_PARAM,
  ORDER_PARAM,
  ASCENDING_ORDER_PARAM
];
