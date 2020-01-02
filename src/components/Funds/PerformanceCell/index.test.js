import React from "react";
import { shallow } from "enzyme";
import PerformanceCell from "./index";

const defaultProps = {
  type: "1y"
};

it("should always render as a span", () => {
  expect(shallow(<PerformanceCell {...defaultProps} />).type()).toEqual("span");
});

it("should have style for given value", () => {
  const wrapper = shallow(<PerformanceCell {...defaultProps} value={10} />);

  expect(wrapper.props().className.split(" ").length).toEqual(3);
});

it("should have style for negative value", () => {
  const wrapper = shallow(<PerformanceCell {...defaultProps} value={-10} />);

  expect(wrapper.props().className.split(" ").length).toEqual(3);
});
