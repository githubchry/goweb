#include "person_detection_wrapper.h"
#include "algorithm.pb.h"

int person_detection_wrapper(char *input, int inlen, char *output, int outlen)
{
	logics::AlgorithmInput imgs;
	std::string img;
	if (imgs.ParseFromArray(input, inlen))
	{
		for (int i = 0; i < imgs.image_size(); i++)
		{
			img = imgs.image(i);
		}
	}

	return img.length();
}