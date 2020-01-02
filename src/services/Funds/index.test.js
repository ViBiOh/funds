import ConfigService from "services/Config";
import FundsService from "./index";

jest.mock("services/Config");

it("should fetch data for funds", () => {
  const getMock = jest.fn().mockReturnValue({
    results: []
  });

  ConfigService.getClient.mockReturnValue({
    get: getMock
  });

  return FundsService.getFunds().then(() => {
    expect(getMock).toHaveBeenCalledWith("/list");
  });
});
