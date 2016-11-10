import React, { Component } from 'react';
import FaClose from 'react-icons/lib/fa/close';
import FaFilter from 'react-icons/lib/fa/filter';
import FaSortAmountAsc from 'react-icons/lib/fa/sort-amount-asc';
import FaSortAmountDesc from 'react-icons/lib/fa/sort-amount-desc';
import { buildFullTextRegex } from '../FullTextSearch';
import Throbber from '../Throbber/Throbber';
import MorningStarService from './MorningStarService';
import MorningStarRow from './MorningStarRow';
import style from './MorningStar.css';

const morningStarIdList = [
  'F00000VSR5', 'F0GBR05V29', 'F00000OVC2', 'F00000VRCL', 'F0GBR05SSA', 'F00000VU45', 'F0GBR04SNN',
  'F000000OZU', 'F00000GUW3', 'F0GBR04KF9', 'F00000OLXV', 'F0000003IP', 'F0GBR04IPB', '0P00006PX0',
  'F0GBR04AMP', 'F00000406S', 'F0GBR04M06', 'F00000JQYK', 'F0GBR04M0L', 'F0GBR05SYI', 'F0GBR04ARL',
  'F0000028FE', 'F0GBR04D1E', 'F0GBR04D20', 'F00000O831', 'F0GBR0636A', 'F0GBR04GA7', 'F0GBR04D14',
  'F0GBR04RHN', 'F00000SF8W', 'F0GBR04FFI', 'F0GBR06368', 'F0GBR04ONO', 'F00000X3BN', 'F0GBR06QR9',
  'F0GBR04I1D', 'F00000VQRW', 'F0GBR04ARH', 'F0GBR06CQU', 'F00000MIEA', 'F0GBR04P4B', 'F00000OKBC',
  'F0GBR04SJG', 'F00000MK7J', 'F0GBR04ASD', 'F0GBR04EC4', 'F0GBR04GQY', 'F0GBR05S8J', 'F0GBR04D1I',
  'F0GBR04KF7', 'F0GBR04RVD', 'F0GBR04DTD', 'F00000SFZQ', 'F0GBR067SB', 'F0GBR04VGJ', 'F0GBR04DTB',
  'F0GBR06OOW', 'F0GBR04ARO', 'F0GBR064XS', 'F00000PXI9', 'F0GBR04JJW', 'F0GBR06KVQ', 'F00000PXEX',
  'F0GBR04CMA', 'F000001EWH', 'F00000JQYQ', 'F00000QLVQ', 'F00000PBTW', 'F0GBR04ONG', 'F0GBR04FBV',
  'F0GBR04ELK', 'F0GBR060JU', 'F0GBR05V5V', 'F000002Q9K', 'F0GBR04VZA', 'F0GBR04ES1', 'F0GBR04D0J',
  'F0GBR04J1N', 'F0GBR04QCJ', 'F000001EWI', 'F0GBR05KO7', 'F0GBR04MWB', 'F0GBR04D1K', 'F0GBR04NA8',
  'F0GBR04NSD', 'F0GBR06PWA', 'F0000000R2', 'F00000JQZ9', 'FOGBR05KMC', 'F0000025P4', 'F0GBR04F8J',
  'F0GBR04BR2', 'F0GBR04OYQ', 'F00000269A', 'F0GBR04QSR', 'F00000MBFD', 'F00000PXB4', 'F0GBR04QJD',
  'F0GBR066X9', 'F0GBR04VMZ', 'F0GBR04FFK', 'F0GBR04AW1', 'F0GBR04HFS', 'F00000THIX', 'F00000S8F1',
  'F00000IS7G', 'F0GBR04SJ0', 'F0GBR04GCU', 'F0GBR04EQE', 'F0GBR05XAB', 'F0GBR04HTO', 'F000000224',
  'F0GBR04W0Z', 'F00000WC76', 'F0GBR055WP', 'F000001AKN', 'F0GBR06D73', 'F0GBR04QII', 'F0GBR04FUQ',
  'FOGBR05JQG', 'F0GBR05XS7', 'F000001T3K', 'F0GBR060ZD', '0P00000JHG', 'F00000SXCI', 'F0GBR05VWK',
  'F00000VXLO', '0P0000941D', 'F0GBR04D1Q', 'F0GBR04MAJ', 'F0GBR04M1Q', 'F0GBR065P8', 'F0GBR05V32',
  'F0GBR06PWC', 'F0GBR04TXT', 'F0GBR04ETZ', 'F0GBR04G2C', 'F00000PXGS', 'FOGBR05JGX', 'F00000WL9G',
  'F0GBR04ASN', 'F0GBR04ARK', 'F0GBR04EI8', 'F000001LNO', 'F0GBR04SEW', 'F0GBR04H58', 'F0GBR06IG2',
  'F0GBR04H31', 'F0GBR04OE5', 'F00000MC3I', 'F0GBR04I5R', 'F0GBR04AXR', 'F0GBR04I54', 'F00000MEBD',
  'F0GBR06LXM', 'F0GBR04QDR', 'F0GBR04HS3', 'F0GBR04G0K', 'F00000S8W7', 'F0GBR04DB3', 'F00000PXD1',
  'F000002KUP', 'F0000007LD', 'F0GBR04ONS', 'F0GBR04RDN', 'F00000PXGZ', 'F000002180', 'F0GBR06T6F',
  'F0GBR04K1R', 'F0GBR068MN', 'F0GBR04RBR', 'F00000V7KE', 'F000001FSI', 'F0GBR04BYP', 'F0GBR04SN7',
  'F0GBR04K9L', 'F0GBR04QCE', 'F0GBR04NAJ', 'F0GBR06BEH', 'F0GBR04BD2', 'F0GBR04F8W', 'F0GBR05SD4',
  'F0GBR04T6H', 'F0GBR04RDG', 'F0GBR065A5', 'F0000043PS', 'F0GBR04EP6', 'F0GBR06IGI', 'F0000021A6',
  'F0GBR06T9L', 'F0GBR04PD2', 'F0GBR04RFU', 'F00000JQZD', 'F00000PXI6', 'F0GBR060RC', 'F0GBR054UL',
  'F0GBR04F94', 'F000000AVO', '0P0000JT24', 'F000002IPX', 'F0GBR06T9O', 'F0GBR04QX4', 'F00000V1UU',
  'F0GBR04IC6', 'F0GBR04D1Y', 'F0GBR04EAW', 'F00000PBVV', 'F0GBR04VCA', 'F0GBR04SNQ', 'F0GBR04KDU',
  'F0GBR04GUI', 'F0GBR06TFP', 'F00000NDR0', 'F0GBR04QCM', 'F0GBR04R2B', 'F0GBR04QTF', 'F0GBR04QDN',
  '0P00000HUX', 'F00000UIMH', 'F0GBR06EC8', 'F000003WD6', 'F0GBR04RKB', 'F0GBR0692A', 'F0GBR04QTQ',
  'F0GBR06X7H', 'F0GBR04AMK', 'F0GBR04OW5', 'F0GBR04GPT', 'F0GBR04BC7', 'F0GBR04F90', 'F0GBR061JX',
  'F0GBR06I4H', 'F00000N9TA', 'F0GBR04F92', 'F0GBR04G92', 'F0GBR04TXJ', 'F0GBR054GE', 'F0GBR04S23',
  'F00000V7KK', 'F0GBR04AS2', 'F0GBR04H3C', 'F0GBR04VG4', 'F0GBR06TGV', '0P0000OQNK', 'F0GBR04R60',
  'F0GBR04FTK', 'F000000ES6', 'F0GBR04QEH', 'F0GBR04QV0', 'F0GBR04IC2', 'F0GBR04F5B', 'F0GBR04EEL',
  'F00000WB2K', '0P00001BUG', 'F0GBR04CBZ', 'F0000002EC', 'F0GBR06IBC', 'F0GBR04QKC', 'F0GBR04D0P',
  'F0GBR04R2V', 'F0GBR04R9T', 'F00000WBSB', 'F0GBR04G5D', 'F0GBR05SGD', 'F0GBR04H5A', 'F00000WL57',
  'F0GBR04CBH', 'F0GBR06MZ9', 'F0GBR04QOM', 'F00000NY9A', 'F0GBR04G59', 'F00000PATC', 'F0GBR04RBZ',
  'F0GBR04AR1', 'F0GBR04G3N', 'F00000W9MR', 'F0GBR04P0V', 'F0GBR04G25', 'F0GBR04LMO', 'F0GBR04VAN',
  'F0GBR04QX6', 'F0GBR04RHP', 'F0GBR04N7E', 'F0GBR04QL4', 'F0000025F0', 'F0000023TI', 'F0GBR04IRU',
  'F0GBR04RVF', 'F0GBR04QTR', 'F0GBR056R2', 'F0GBR066P1', 'F0GBR04NAD', 'F000003XBR', 'F00000MJ9P',
  'F0000007LH', 'F0GBR04NSH', 'F0GBR05U0S', 'F0GBR04SNI', 'F000005OFD', 'F00000LMV2', 'F0GBR04OA5',
  'F0GBR04HKT', 'F0GBR04ZW9', 'F0GBR04BG8', 'F00000J3T6', '0P00005ZG6', 'F0GBR05KM2', 'F00000TP9D',
  'F0GBR04HU5', 'F000001V1S', 'F0GBR04QIL', 'F00000202C', 'F0GBR04FTI', 'F00000H6EA', 'F0GBR04DMP',
  'F00000UDIE', 'F0GBR04SGP', 'F0GBR04D26', 'F0000026XE', 'F00000V7KH', '0P0000OP65', '0P0000191V',
  'F0GBR04BWM', 'F0000026T3', 'F0GBR04TBJ', 'F000000G95', 'F0GBR04KF3', 'F0GBR04D0L', 'F0GBR04HRH',
  'F0GBR04F8U', 'F0GBR04M1F', 'F00000GXW5', 'F0GBR04QMW', 'F0GBR04MAD', 'F0GBR04D0H', 'F0GBR04EVF',
  'F0GBR04CBX', 'F00000UO20', 'F000000G8T', 'F000002RAX', 'F0GBR04FVQ', 'F0GBR05Z9V', 'F0GBR04I6R',
  'F0GBR04F8Y', 'F0GBR04VGP', 'F000002N3W', 'F0GBR04JN6', 'F0000003C1', 'F00000WDEN', 'F0GBR0697H',
  'F00000OB43', 'F000002LA1', 'F0GBR04RBU', 'F000002RAY', 'F0GBR04IT2', 'F0GBR04R5Y', 'F0GBR05T1R',
  'F0000025EJ', 'F0GBR04G69', 'F00000MJWQ', 'F0GBR04D1G', 'FOGBR05KOD', 'F0GBR05RX0', 'F0GBR04QBU',
  'F0GBR04R6G', 'F0GBR04QUS', 'F0GBR067EX', 'F000000PD0', 'F0GBR04G8M', 'F000002D7K', 'F0GBR04SO4',
  'F0GBR04VUL', 'F0GBR04VT4', 'F0GBR04BDA', 'F0GBR04SW9', 'F0GBR04T5V', 'F0GBR04F3U', 'F0GBR04EF5',
  'F0GBR05SSC', '0P00000MRP', 'F0GBR04AWP', 'F0GBR04N26', 'F00000WL99', 'F0GBR04REK', 'F0GBR04H39',
  'F00000WL6B', 'F0GBR04M00', 'F00000QGES', 'F0GBR05SCF', 'F0GBR04G8K', 'F0GBR04FCH', 'F000003YDQ',
  'F0GBR04SNC', 'FOGBR05KLS', 'F0GBR04D2C', 'F0GBR04D0V', 'F0GBR04STB', 'F0GBR06T9G', 'F0GBR04YW1',
  'F00000NUAT', 'F0GBR04RL1', 'F0GBR06NZD', 'F00000QB80', 'F0GBR04G0A', 'F0GBR04G88', 'F0GBR04FGC',
  'F0GBR04QL5', 'F0GBR05TLP', 'F0GBR064H6', 'F000002ET8', 'F0GBR04K8F', 'F0GBR04R2J', 'F0GBR06PNW',
  'F0GBR04M2G', '0P00000MRI', 'FOGBR05JGL', 'F00000SFI4', 'F0GBR04IU2', 'F0GBR04QX8', 'F00000WL6H',
  'F00000PHWD', 'F0GBR04MYA', 'F0GBR06MIL', 'F000000LUN', 'F0GBR04H5P', 'F0GBR04D2G', 'F0GBR04LMM',
  'F0GBR05V8O', 'F0GBR04IWV', 'F0GBR06083', 'F00000PXES', 'F0GBR04QV8', 'F0GBR04KPS', 'F0GBR04M9Y',
  'F0GBR04QIY', 'F0GBR04RZG', 'F00000WL6T', 'F0GBR06OFA', 'F00000MOTL', 'F0GBR04CBL', 'F00000J81V',
  'F00000MKFB', 'F0GBR06BNT', 'F0GBR04AFY', 'F0GBR04R8S', 'F0GBR0642C', 'F0GBR05YT7', 'F0GBR05U73',
  'F0GBR06HDW', 'F0GBR04R14', 'F0GBR04AG0', 'F0GBR04LMY', 'F0GBR05TGT', 'F0GBR04R0R', 'F0GBR04SNE',
  'F0GBR04JY3', 'F000000ESQ', 'F000000ET1', 'F0GBR04D0R', 'F0GBR0580R', 'F00000UILU', 'F0GBR04ODP',
  'F0GBR04EF8', 'F0GBR04SNR', 'F0GBR05VWR', 'F000000F7B', 'F0GBR04EN1', 'F0GBR04FX5', 'F0GBR04DM7',
  'F0GBR04AXJ', 'F0GBR04QCF', 'F0GBR05VVH', 'F0GBR069YZ', 'F0GBR04N1E', 'F0GBR04RGK', 'F0GBR04QWF',
  'F0GBR04N04', 'F0GBR04QUK', 'F0GBR04K6T', 'F0GBR04BXB', 'F0GBR04NPY', 'F0GBR04PF1', 'F0GBR04EUJ',
  'F00000T6G4', 'FOGBR05K58', '0P0000M6N2', 'F000000GT1', 'F0GBR04DZL', 'F000000L7N', 'F0GBR04D1C',
  'F0GBR04SIW', 'F00000VQIN', 'F0GBR04QVL', 'F0GBR04M1S', 'F0GBR04D0T', 'F00000H4HR', 'F0GBR04JGM',
  'F0GBR04GTL', 'F0GBR04REQ', 'F0GBR04QK4', '0P0000G6K3', 'F0GBR04D11', 'F0GBR04QJ0', 'F0GBR04EB8',
  'F0GBR04QOY', 'F0GBR04D22', 'F00000H0V9', 'F0GBR05Z1J', 'F0GBR04QCK', 'F0GBR04D0X', 'F0GBR04AO3',
  'F0GBR04EF6', 'F0GBR04HZ7', 'F0GBR04MZY', 'F0000000DH', 'FOGBR05JLL', 'F0GBR05V7K', 'F00000J58Y',
  'F0GBR04VH3', '0P00009419', 'FOGBR05K54', 'F00000VXRH', 'F00000PXH7', 'F0GBR04D18', 'F0GBR04QFV',
  'F0GBR068QC', 'F000000QY1', 'F0GBR06HJU', 'F0GBR04E8X', 'F0GBR063QR', 'F0GBR04BHX', 'F0GBR04JNP',
  'F0GBR0636J',
];

const FETCH_SIZE = 25;

const COLUMNS = {
  isin: {
    label: 'ISIN',
    sortable: false,
  },
  label: {
    label: 'Libellé',
    sortable: true,
  },
  category: {
    label: 'Catégorie',
    sortable: true,
  },
  rating: {
    label: 'Note',
    sortable: true,
  },
  '1m': {
    label: '1 mois',
    sortable: true,
  },
  '3m': {
    label: '3 mois',
    sortable: true,
  },
  '6m': {
    label: '6 mois',
    sortable: true,
  },
  '1y': {
    label: '1 an',
    sortable: true,
  },
  v1y: {
    label: 'Volatilité',
    sortable: true,
  },
  score: {
    label: 'Score',
    sortable: true,
  },
};

export default class MorningStarList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      performances: [],
      displayed: [],
      order: {
        key: '',
        descending: true,
      },
      filters: {
        text: '',
      },
      orderDisplayed: false,
    };

    this.fetchPerformances = this.fetchPerformances.bind(this);
    this.fetchPerformance = this.fetchPerformance.bind(this);

    this.filterBy = this.filterBy.bind(this);
    this.orderBy = this.orderBy.bind(this);
    this.onOrderClick = this.onOrderClick.bind(this);
    this.reverseOrder = this.reverseOrder.bind(this);

    this.updateDataPresentation = this.updateDataPresentation.bind(this);
    this.renderRow = this.renderRow.bind(this);
    this.renderOrder = this.renderOrder.bind(this);
    this.renderSearch = this.renderSearch.bind(this);
  }

  componentDidMount() {
    for (let i = 0, size = morningStarIdList.length; i < size; i += FETCH_SIZE) {
      this.fetchPerformances(morningStarIdList.slice(i, i + FETCH_SIZE));
    }
  }

  onOrderClick() {
    this.setState({ orderDisplayed: !this.state.orderDisplayed });
  }

  fetchPerformances(ids) {
    return MorningStarService.getPerformances(ids)
      .then((performances) => {
        const results = performances.results.filter(performance => performance.id);
        this.setState({
          performances: [...this.state.performances, ...results],
        }, this.updateDataPresentation);

        return performances;
      });
  }

  fetchPerformance(id) {
    return MorningStarService.getPerformance(id)
      .then((performance) => {
        this.setState({
          performances: [...this.state.performances, performance],
        }, this.updateDataPresentation);

        return performance;
      });
  }

  filterBy(filterName, value) {
    const filter = {};
    filter[filterName] = value;

    this.setState({
      filters: Object.assign(this.state.filters, filter),
    }, this.updateDataPresentation);
  }

  orderBy(order) {
    this.setState({
      order: Object.assign(this.state.order, { key: order, descending: true }),
    }, this.updateDataPresentation);

    if (order) {
      this.onOrderClick();
    }
  }

  reverseOrder() {
    this.setState({
      order: Object.assign(this.state.order, { descending: !this.state.order.descending }),
    }, this.updateDataPresentation);
  }

  updateDataPresentation() {
    clearTimeout(this.timeout);
    this.timeout = setTimeout(() => {
      let displayed = this.state.performances.slice();

      const categoryFilter = this.state.filters.category;
      if (categoryFilter) {
        displayed = displayed.filter(({ category }) => category === categoryFilter);
      }

      const ratingFilter = this.state.filters.rating;
      if (ratingFilter) {
        displayed = displayed.filter(({ rating }) => rating === ratingFilter);
      }

      if (this.state.filters.text) {
        const regex = buildFullTextRegex(this.state.filters.text);
        displayed = displayed.filter(({ label }) => regex.test(label));
      }

      if (this.state.order.key) {
        const orderKey = this.state.order.key;
        const compareMultiplier = this.state.order.descending ? -1 : 1;

        displayed = displayed.sort((o1, o2) => {
          if (!o1 || typeof o1[orderKey] === 'undefined') {
            return -1 * compareMultiplier;
          } else if (!o2 || typeof o2[orderKey] === 'undefined') {
            return 1 * compareMultiplier;
          } else if (o1[orderKey] < o2[orderKey]) {
            return -1 * compareMultiplier;
          } else if (o1[orderKey] > o2[orderKey]) {
            return 1 * compareMultiplier;
          }
          return 0;
        });
      }

      this.setState({
        displayed,
      });
    }, 400);
  }

  renderRow() {
    return this.state.displayed.map(performance => (
      <MorningStarRow
        key={performance.id}
        performance={performance}
        filterBy={this.filterBy}
      />
    ));
  }

  renderOrder() {
    const orderColumns = Object.keys(COLUMNS)
      .filter(column => COLUMNS[column].sortable)
      .map(key => (
        <li key={key}><button onClick={() => this.orderBy(key)}>{COLUMNS[key].label}</button></li>
      ));

    return (
      <span className={style.icon}>
        <FaSortAmountDesc
          className={this.state.orderDisplayed ? style.active : ''}
          onClick={this.onOrderClick}
        />
        <ol className={this.state.orderDisplayed ? style.displayed : style.hidden}>
          {orderColumns}
        </ol>
      </span>
    );
  }

  renderSearch() {
    return (
      <span className={style.search}>
        <input
          type="text"
          placeholder="Filter on label"
          value={this.state.filters.text}
          onChange={e => this.filterBy('text', e.target.value)}
        />
        <span className={style.count}>
          {this.state.displayed.length} / {morningStarIdList.length}
        </span>
        {this.renderOrder()}
      </span>
    );
  }

  render() {
    if (this.state.performances.length === 0) {
      return <Throbber label="Chargement des données" />;
    }

    const header = Object.keys(COLUMNS).reduce((previous, current) => {
      previous[current] = COLUMNS[current].label; // eslint-disable-line no-param-reassign
      return previous;
    }, {});

    return (
      <span>
        <div key="search" className={style.list}>
          {this.renderSearch()}
          {
            this.state.filters.category && (
              <span className={style.dataModifier}>
                <span className={style.icon}>
                  <FaFilter />
                </span>
                <span><em> Catégorie</em> = </span>
                {this.state.filters.category}
                <button onClick={() => this.filterBy('category', '')} className={style.icon}>
                  <FaClose />
                </button>
              </span>
            )
          }
          {
            this.state.filters.rating && (
              <span className={style.dataModifier}>
                <span className={style.icon}>
                  <FaFilter />
                </span>
                <span><em> Note</em> = </span>
                {this.state.filters.rating}
                <button onClick={() => this.filterBy('rating', '')} className={style.icon}>
                  <FaClose />
                </button>
              </span>
            )
          }
          {
            this.state.order.key && (
              <span className={style.dataModifier}>
                <button onClick={this.reverseOrder} className={style.icon}>
                  {
                    this.state.order.descending ? <FaSortAmountDesc /> : <FaSortAmountAsc />
                  }
                </button>
                &nbsp;{COLUMNS[this.state.order.key].label}
                <button onClick={() => this.orderBy('')} className={style.icon}>
                  <FaClose />
                </button>
              </span>
            )
          }
        </div>
        <div key="list" className={style.list}>
          <MorningStarRow key={'header'} performance={header} />
          {this.renderRow()}
        </div>
      </span>
    );
  }
}
