# yardsard

A reimplementation of the Yard Sale Model in golang. The original version of this
simulation can be found at [https://physics.umd.edu/hep/drew/math_general/yard_sale.html].
The Yard Sale model is a simplified demonstration of how small differences in
valuation of goods can lead to the centralization of wealth over time, even when
the only factor contributing to the difference is random chance. To quote the
original site:

> To understand the mathematics of wealth inequality, start with what is known
as the "Yard Sale Model", based on a rather well defined economic system that
describes the exchange of goods among people in a fixed population.
> In this model, there is no currency, but there is lots of trading among pairs
of people.
> Imagine what would happen if every trade is fair (trading goods of equal value).
> Then, over time, everyone's net wealth will stay constant, and the net wealth
of the community stays constant.
> Now imagine what would happen if the trades are sometimes not completely equal
(by some measure).
> What happens over time is pretty interesting.

This project is in its early stages, and deviates from the original model in a
few distinct ways:

1. Individuals within the popluation (referred to as "agents") are not given the
same starting wealth, but instead a random amount of money between $1 and the maximum
starting wealth, to simulate the existing disparity between rich and poor.
2. If the total allocated wealth after the first step is not equal to the total
available wealth if all agents were given the same starting value, money is randomly
assigned until the two values are equal.
3. The simulation stops when one agent has accumulated 95% of the total available
wealth, rather than running indefinitely.
4. Currently, only the distribution of wealth across all agents is plotted, with
no histogram to display the wealth of an individual over time.
